package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

// nominatimResponse models the fields we read from a Nominatim reverse
// geocoding response.
type nominatimResponse struct {
	DisplayName string `json:"display_name"`
	Address     struct {
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		County      string `json:"county"`
		State       string `json:"state"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

// Geocoder reverse-geocodes GPS coordinates into a human-readable location
// using a configurable Nominatim-compatible endpoint. Results are cached to
// avoid repeated lookups for the same coordinates.
type Geocoder struct {
	endpoint   string
	userAgent  string
	client     *http.Client
	cache      port.CacheService
}

// NewGeocoder creates a Geocoder. The endpoint must be a Nominatim-compatible
// server (e.g. https://nominatim.openstreetmap.org or a self-hosted instance).
func NewGeocoder(config appconfig.AppConfig, cache port.CacheService) *Geocoder {
	endpoint := strings.TrimRight(config.GeocodeEndpoint, "/")
	if endpoint == "" {
		endpoint = "https://nominatim.openstreetmap.org"
	}
	return &Geocoder{
		endpoint:  endpoint,
		userAgent: "photobox/" + "1.0",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		cache: cache,
	}
}

// Reverse returns a short human-readable location (city, state) for the given
// coordinates. Returns ("", nil) when the endpoint is unreachable or the
// coordinates are empty, so failures degrade gracefully.
func (g *Geocoder) Reverse(ctx context.Context, latitude, longitude float64) (string, error) {
	if latitude == 0 && longitude == 0 {
		return "", nil
	}

	cacheKey := fmt.Sprintf("geocode:%.6f:%.6f", latitude, longitude)
	if g.cache != nil {
		if cached, err := g.cache.Get(cacheKey); err == nil && cached != "" {
			return cached, nil
		}
	}

	u, err := url.Parse(fmt.Sprintf("%s/reverse", g.endpoint))
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("lat", fmt.Sprintf("%.6f", latitude))
	q.Set("lon", fmt.Sprintf("%.6f", longitude))
	q.Set("format", "jsonv2")
	q.Set("zoom", "12")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", g.userAgent)

	resp, err := g.client.Do(req)
	if err != nil {
		slog.Warn("reverse geocoding request failed", "error", err)
		return "", nil
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("reverse geocoding returned non-200", "status", resp.StatusCode)
		return "", nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}

	var data nominatimResponse
	if err := json.Unmarshal(body, &data); err != nil {
		slog.Warn("reverse geocoding response parse failed", "error", err)
		return "", nil
	}

	location := data.shortName()
	if location != "" && g.cache != nil {
		_ = g.cache.Set(cacheKey, location, 30*24*time.Hour)
	}
	return location, nil
}

// shortName composes "City, Country" (or town/village/county fallbacks) from
// the address parts.
func (n *nominatimResponse) shortName() string {
	parts := []string{n.Address.City, n.Address.Town, n.Address.Village, n.Address.County}
	var locality string
	for _, p := range parts {
		if p != "" {
			locality = p
			break
		}
	}
	if locality != "" && n.Address.Country != "" {
		return fmt.Sprintf("%s, %s", locality, n.Address.Country)
	}
	if locality != "" {
		return locality
	}
	return n.Address.Country
}
