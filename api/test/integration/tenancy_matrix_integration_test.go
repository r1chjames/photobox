//go:build integration

package integration

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/repository"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"gitlab.com/r1chjames/photobox/api/test/testutil"
	"gorm.io/datatypes"
)

// TestCrossTenantContentMatrix_Integration is the definition-of-done gate for
// the Phase 2 query-scoping sweep (issue #74 §7): two tenants A and B, where
// A must never read or mutate B's content on any endpoint.
//
// It is gated behind TENANCY_MATRIX=1 because it is expected to FAIL until the
// repository queries are scoped by workspace — running it by default would
// leave CI red for a known, tracked gap. Run it to enumerate exactly which
// endpoints leak:
//
//	TENANCY_MATRIX=1 go test -tags=integration -run TestCrossTenantContentMatrix ./test/integration/ -v
//
// When the sweep lands, drop the gate and make this a normal test.
func TestCrossTenantContentMatrix_Integration(t *testing.T) {
	if os.Getenv("TENANCY_MATRIX") != "1" {
		t.Skip("set TENANCY_MATRIX=1 to run the cross-tenant content matrix (fails until Phase 2 query scoping lands)")
	}

	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	photoRepo := repository.NewPhotoRepository(env)
	albumRepo := repository.NewAlbumRepository(env)
	shareRepo := repository.NewShareRepository(env)

	// Two isolated tenants, each with a personal workspace (§8).
	alice, aliceWS, err := userService.RegisterWithPersonalWorkspace(&domain.User{
		Username: "matrix-alice", Email: "matrix-alice@example.com", Password: "password123",
	}, "Alice Photos")
	require.NoError(t, err)
	require.NotEmpty(t, alice.ID)
	bob, bobWS, err := userService.RegisterWithPersonalWorkspace(&domain.User{
		Username: "matrix-bob", Email: "matrix-bob@example.com", Password: "password123",
	}, "Bob Photos")
	require.NoError(t, err)
	require.NotEqual(t, aliceWS.ID, bobWS.ID, "tenants must be isolated")

	// Bob's content, explicitly in Bob's workspace.
	bobAlbum, err := albumRepo.CreateAlbum("Bob Album")
	require.NoError(t, err)
	require.NoError(t, env.Db.Model(&domain.Album{}).Where("id = ?", bobAlbum.ID).
		Update("workspace_id", bobWS.ID).Error)

	bobPhoto := domain.Photo{
		ID: "bob-photo-1", Name: "bob-secret.jpg", AlbumId: bobAlbum.ID,
		FilesystemPath: "/photos/bob/bob-secret.jpg", MediaType: "image",
		CreatedEpoch: time.Now().UnixMilli(), Metadata: datatypes.JSON([]byte("{}")),
	}
	require.NoError(t, photoRepo.CreatePhotoInfo(bobPhoto))
	require.NoError(t, env.Db.Model(&domain.Photo{}).Where("id = ?", bobPhoto.ID).
		Update("workspace_id", bobWS.ID).Error)

	bobShare := &domain.SharedLink{Token: "bob-share-token", ResourceType: "photo", ResourceId: bobPhoto.ID, CreatedBy: bob.ID}
	require.NoError(t, shareRepo.CreateShare(bobShare))
	require.NoError(t, env.Db.Model(&domain.SharedLink{}).Where("token = ?", bobShare.Token).
		Update("workspace_id", bobWS.ID).Error)

	// Actor is Alice, acting in her own workspace.
	config := testutil.TestDBConfig()
	config.Token = "0123456789abcdef0123456789abcdef"
	config.TokenDuration = time.Hour
	config.CorsAllowedOrigins = []string{"http://localhost"}
	config.ApiBasePath = "/api"
	router := buildAuthzRouter(t, env, &config)
	aliceToken := loginAs(t, router, "matrix-alice", "password123")
	assert.NotEmpty(t, aliceToken)

	authHeaders := func() map[string]string {
		return map[string]string{"Authorization": "Bearer " + aliceToken, "X-Workspace-ID": aliceWS.ID}
	}

	// Every case: Alice targets Bob's resource. A 2xx on a single-resource
	// endpoint is a breach; a 2xx on a list/search endpoint is a breach when
	// the body contains Bob's identifiers.
	cases := []struct {
		name           string
		method         string
		path           string
		body           any
		breachWhen2xx  bool
		mustNotContain []string
	}{
		{"read photo detail", http.MethodGet, "/api/photo/info/" + bobPhoto.ID, nil, true, nil},
		{"read photo binary", http.MethodGet, "/api/photo/bin/" + bobPhoto.ID, nil, true, nil},
		{"read photo thumbnail", http.MethodGet, "/api/photo/thumbnail/" + bobPhoto.ID + "?size=m", nil, true, nil},
		{"delete photo", http.MethodDelete, "/api/photos/" + bobPhoto.ID, nil, true, nil},
		{"favorite photo", http.MethodPatch, "/api/photos/" + bobPhoto.ID + "/favorite", map[string]any{"favorite": true}, true, nil},
		{"update photo tags", http.MethodPatch, "/api/photos/" + bobPhoto.ID + "/tags", map[string]any{"tags": "hacked"}, true, nil},
		{"read album", http.MethodGet, "/api/album/" + bobAlbum.ID, nil, true, nil},
		{"delete album", http.MethodDelete, "/api/albums/" + bobAlbum.ID, nil, true, nil},
		{"list photos", http.MethodGet, "/api/photos?limit=100", nil, false, []string{bobPhoto.ID, "bob-secret"}},
		{"list albums", http.MethodGet, "/api/albums?limit=100", nil, false, []string{bobAlbum.ID, "Bob Album"}},
		{"search", http.MethodGet, "/api/search?q=secret", nil, false, []string{bobPhoto.ID, "bob-secret"}},
		{"timeline", http.MethodGet, "/api/photos/timeline", nil, false, nil},
		{"all tags", http.MethodGet, "/api/photos/tags", nil, false, nil},
		{"duplicates", http.MethodGet, "/api/photos/duplicates", nil, false, []string{bobPhoto.ID, "bob-secret"}},
		{"trash list", http.MethodGet, "/api/photos/trash?limit=100", nil, false, []string{bobPhoto.ID, "bob-secret"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := doRequest(router, tc.method, tc.path, tc.body, authHeaders())
			is2xx := w.Code >= 200 && w.Code < 300
			if !is2xx {
				return
			}
			if tc.breachWhen2xx {
				t.Errorf("LEAK: Alice %s %s -> %d returned another tenant's resource\nbody: %s",
					tc.method, tc.path, w.Code, truncate(w.Body.String(), 200))
				return
			}
			for _, marker := range tc.mustNotContain {
				if strings.Contains(w.Body.String(), marker) {
					t.Errorf("LEAK: Alice %s %s -> %d response contains another tenant's data (%q)\nbody: %s",
						tc.method, tc.path, w.Code, marker, truncate(w.Body.String(), 300))
					return
				}
			}
			// No marker observed is NOT proof of isolation — the endpoint may
			// simply be empty or filter by other means. Report it so a PASS
			// here is never read as verified isolation.
			t.Logf("UNVERIFIED: Alice %s %s -> %d returned no other-tenant marker; confirm the query is workspace-scoped",
				tc.method, tc.path, w.Code)
		})
	}
}

// truncate shortens a response body for readable failure output.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
