//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httpHandler "gitlab.com/r1chjames/photobox/api/internal/adapter/handler/http"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/repository"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"gitlab.com/r1chjames/photobox/api/test/testutil"
)

// doRequest issues a JSON request through the real router with optional
// headers and returns the recorder.
func doRequest(router *httpHandler.Router, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)
	return w
}

// TestWorkspaceLifecycle_Integration covers workspace CRUD, owner-on-create,
// membership management, and the last-owner rule (issue #74).
func TestWorkspaceLifecycle_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	wsRepo := repository.NewWorkspaceRepository(env)
	wsService := service.NewWorkspaceService(wsRepo)

	owner, err := userService.CreateUser(&domain.User{
		Username: "ws-owner", Email: "ws-owner@example.com", Password: "password123",
		Role: domain.CONTRIBUTOR, Approved: true,
	})
	require.NoError(t, err)

	// Create: the creator becomes owner.
	ws, err := wsService.CreateWorkspace("Holiday Photos", "", owner.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, ws.ID)
	assert.Equal(t, "holiday-photos", ws.Slug)

	m, err := wsRepo.GetMembership(ws.ID, owner.ID)
	require.NoError(t, err)
	require.NotNil(t, m)
	assert.Equal(t, domain.WorkspaceOwner, m.Role)

	// A second workspace gets a distinct slug when the name collides.
	ws2, err := wsService.CreateWorkspace("Holiday Photos", "holiday-photos-2", owner.ID)
	require.NoError(t, err)
	assert.NotEqual(t, ws.ID, ws2.ID)

	// List returns both created workspaces plus the shared default workspace
	// the admin-created user automatically joined.
	list, err := wsService.ListWorkspaces(owner.ID)
	require.NoError(t, err)
	assert.Len(t, list, 3)
	slugs := map[string]bool{}
	for _, w := range list {
		slugs[w.Slug] = true
	}
	assert.True(t, slugs["holiday-photos"] && slugs["holiday-photos-2"] && slugs["default"])

	// A non-member cannot read the workspace.
	_, err = wsService.GetWorkspace(ws.ID, "someone-else")
	assert.ErrorIs(t, err, domain.ErrForbidden)

	// Add a member, then verify their role and access.
	wsService.AddMember(ws.ID, owner.ID, "member-user", domain.WorkspaceMemberRole)
	mem, err := wsRepo.GetMembership(ws.ID, "member-user")
	require.NoError(t, err)
	require.NotNil(t, mem)
	assert.Equal(t, domain.WorkspaceMemberRole, mem.Role)

	// A plain member cannot manage members (403).
	err = wsService.AddMember(ws.ID, "member-user", "third-user", domain.WorkspaceViewer)
	assert.ErrorIs(t, err, domain.ErrForbidden)

	// Rename requires manage rights; a member is refused.
	_, err = wsService.UpdateWorkspace(ws.ID, "member-user", "Renamed")
	assert.ErrorIs(t, err, domain.ErrForbidden)
	renamed, err := wsService.UpdateWorkspace(ws.ID, owner.ID, "Renamed")
	require.NoError(t, err)
	assert.Equal(t, "Renamed", renamed.Name)

	// Last-owner rule: an owner cannot be removed while they are the only one.
	assert.Error(t, wsService.RemoveMember(ws.ID, owner.ID, owner.ID))

	// With a second owner present, removal is allowed.
	require.NoError(t, wsService.AddMember(ws.ID, owner.ID, "second-owner", domain.WorkspaceOwner))
	assert.NoError(t, wsService.RemoveMember(ws.ID, owner.ID, "second-owner"))

	// Only the owner may delete the workspace.
	assert.ErrorIs(t, wsService.DeleteWorkspace(ws.ID, "member-user"), domain.ErrForbidden)
	assert.NoError(t, wsService.DeleteWorkspace(ws.ID, owner.ID))

	_, err = wsRepo.GetWorkspaceById(ws.ID)
	assert.ErrorIs(t, err, domain.ErrDataNotFound)
}

// TestWorkspaceSignupIsolation_Integration verifies the signup path creates
// an isolated personal workspace per user (issue #74 §8), and that one
// signup's workspace is invisible to another.
func TestWorkspaceSignupIsolation_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	wsRepo := repository.NewWorkspaceRepository(env)
	wsService := service.NewWorkspaceService(wsRepo)

	alice, aliceWS, err := userService.RegisterWithPersonalWorkspace(&domain.User{
		Username: "alice", Email: "alice@example.com", Password: "password123",
	}, "Alice's Photos")
	require.NoError(t, err)

	bob, bobWS, err := userService.RegisterWithPersonalWorkspace(&domain.User{
		Username: "bob", Email: "bob@example.com", Password: "password123",
	}, "Bob's Photos")
	require.NoError(t, err)

	require.NotEqual(t, aliceWS.ID, bobWS.ID, "each signup gets its own workspace")

	// Alice owns exactly her workspace; Bob's is invisible to her.
	aliceWorkspaces, err := wsService.ListWorkspaces(alice.ID)
	require.NoError(t, err)
	assert.Len(t, aliceWorkspaces, 1)
	assert.Equal(t, aliceWS.ID, aliceWorkspaces[0].ID)

	_, err = wsService.GetWorkspace(bobWS.ID, alice.ID)
	assert.ErrorIs(t, err, domain.ErrForbidden, "Alice must not access Bob's workspace")

	// Owner membership is recorded for each signup.
	for _, pair := range []struct{ wsID, userID string }{{aliceWS.ID, alice.ID}, {bobWS.ID, bob.ID}} {
		m, err := wsRepo.GetMembership(pair.wsID, pair.userID)
		require.NoError(t, err)
		require.NotNil(t, m)
		assert.Equal(t, domain.WorkspaceOwner, m.Role)
	}
}

// TestWorkspaceAdminCreatedUser_JoinsDefaultWorkspace_Integration verifies the
// home-instance semantics: users created by an admin belong to the shared
// default workspace (issue #74 D6), rather than getting an isolated one.
func TestWorkspaceAdminCreatedUser_JoinsDefaultWorkspace_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	wsRepo := repository.NewWorkspaceRepository(env)

	u, err := userService.CreateUser(&domain.User{
		Username: "home-user", Email: "home@example.com", Password: "password123",
		Role: domain.CONTRIBUTOR, Approved: true,
	})
	require.NoError(t, err)

	m, err := wsRepo.GetMembership(domain.DefaultWorkspaceID, u.ID)
	require.NoError(t, err)
	require.NotNil(t, m, "admin-created users join the default workspace")
	assert.Equal(t, domain.WorkspaceMemberRole, m.Role)
}

// TestWorkspaceRoutes_Integration exercises the HTTP surface: listing,
// creating, and the auth requirement.
func TestWorkspaceRoutes_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	config := testutil.TestDBConfig()
	config.Token = "0123456789abcdef0123456789abcdef"
	config.TokenDuration = time.Hour
	config.CorsAllowedOrigins = []string{"http://localhost"}
	config.ApiBasePath = "/api"
	router := buildAuthzRouter(t, env, &config)

	// Unauthenticated listing is rejected.
	w := doRequest(router, http.MethodGet, "/api/workspaces", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	_, err := userService.CreateUser(&domain.User{
		Username: "ws-route-user", Email: "ws-route@example.com", Password: "password123",
		Role: domain.CONTRIBUTOR, Approved: true,
	})
	require.NoError(t, err)
	token := loginAs(t, router, "ws-route-user", "password123")

	// Admin-created user auto-joins the default workspace.
	w = doRequest(router, http.MethodGet, "/api/workspaces", nil, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Create a workspace and confirm it appears in the listing.
	w = doRequest(router, http.MethodPost, "/api/workspaces", map[string]any{"name": "Route Test"}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w = doRequest(router, http.MethodGet, "/api/workspaces", nil, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Route Test")
}
