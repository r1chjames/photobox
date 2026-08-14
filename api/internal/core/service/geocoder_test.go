package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
)

// MockCacheServiceLite is a minimal cache stub for geocoder tests.
type MockCacheServiceLite struct {
	mock.Mock
}

func (m *MockCacheServiceLite) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheServiceLite) Set(key string, value string, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheServiceLite) Ping() error {
	return nil
}

func (m *MockCacheServiceLite) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCacheServiceLite) DeletePattern(pattern string) error {
	return nil
}

func (m *MockCacheServiceLite) GetBytes(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCacheServiceLite) SetBytes(key string, value []byte, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

// TestGeocoderReverse tests the Nominatim reverse-geocoding service.
func TestGeocoderReverse(t *testing.T) {
	t.Run("returns short location from nominatim response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/reverse")
			assert.Equal(t, "52.040000", r.URL.Query().Get("lat"))
			_, _ = w.Write([]byte(`{"display_name":"x","address":{"city":"Cambridge","country":"United Kingdom"}}`))
		}))
		defer srv.Close()

		cache := new(MockCacheServiceLite)
		cache.On("Get", "geocode:52.040000:0.094444").Return("", assert.AnError)
		cache.On("Set", "geocode:52.040000:0.094444", "Cambridge, United Kingdom", mock.Anything).Return(nil)

		g := NewGeocoder(appconfig.AppConfig{GeocodeEndpoint: srv.URL}, cache)
		loc, err := g.Reverse(context.Background(), 52.04, 0.094444)

		assert.NoError(t, err)
		assert.Equal(t, "Cambridge, United Kingdom", loc)
		cache.AssertExpectations(t)
	})

	t.Run("uses cache hit", func(t *testing.T) {
		cache := new(MockCacheServiceLite)
		cache.On("Get", "geocode:52.040000:0.094444").Return("Cached Town", nil)

		g := NewGeocoder(appconfig.AppConfig{}, cache)
		loc, err := g.Reverse(context.Background(), 52.04, 0.094444)

		assert.NoError(t, err)
		assert.Equal(t, "Cached Town", loc)
		cache.AssertExpectations(t)
	})

	t.Run("returns empty for zero coordinates", func(t *testing.T) {
		cache := new(MockCacheServiceLite)
		g := NewGeocoder(appconfig.AppConfig{}, cache)
		loc, err := g.Reverse(context.Background(), 0, 0)

		assert.NoError(t, err)
		assert.Equal(t, "", loc)
	})

	t.Run("degrades gracefully on server error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		cache := new(MockCacheServiceLite)
		cache.On("Get", "geocode:52.040000:0.094444").Return("", assert.AnError)

		g := NewGeocoder(appconfig.AppConfig{GeocodeEndpoint: srv.URL}, cache)
		loc, err := g.Reverse(context.Background(), 52.04, 0.094444)

		assert.NoError(t, err)
		assert.Equal(t, "", loc)
	})
}

// TestNominatimShortName tests the location-name composition.
func TestNominatimShortName(t *testing.T) {
	t.Run("city and country", func(t *testing.T) {
		n := nominatimResponse{}
		n.Address.City = "Cambridge"
		n.Address.Country = "United Kingdom"
		assert.Equal(t, "Cambridge, United Kingdom", n.shortName())
	})

	t.Run("town fallback", func(t *testing.T) {
		n := nominatimResponse{}
		n.Address.Town = "Grantchester"
		n.Address.Country = "United Kingdom"
		assert.Equal(t, "Grantchester, United Kingdom", n.shortName())
	})

	t.Run("country only", func(t *testing.T) {
		n := nominatimResponse{}
		n.Address.Country = "France"
		assert.Equal(t, "France", n.shortName())
	})
}
