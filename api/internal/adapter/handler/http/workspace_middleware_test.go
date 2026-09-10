package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// fakeWorkspaceService is a minimal WorkspaceService for middleware tests.
type fakeWorkspaceService struct {
	// memberships maps "workspaceID|userID" -> role.
	memberships map[string]domain.WorkspaceRole
	// workspaces maps userID -> their workspaces.
	workspaces map[string][]domain.Workspace
}

func newFakeWorkspaceService() *fakeWorkspaceService {
	return &fakeWorkspaceService{
		memberships: map[string]domain.WorkspaceRole{},
		workspaces:  map[string][]domain.Workspace{},
	}
}

func (f *fakeWorkspaceService) add(userID, wsID string, role domain.WorkspaceRole) {
	f.memberships[wsID+"|"+userID] = role
	f.workspaces[userID] = append(f.workspaces[userID], domain.Workspace{ID: wsID})
}

func (f *fakeWorkspaceService) CreateWorkspace(string, string, string) (*domain.Workspace, error) {
	panic("not used")
}

func (f *fakeWorkspaceService) GetWorkspace(wsID, userID string) (*domain.Workspace, error) {
	if _, ok := f.memberships[wsID+"|"+userID]; !ok {
		return nil, domain.ErrForbidden
	}
	return &domain.Workspace{ID: wsID}, nil
}

func (f *fakeWorkspaceService) ListWorkspaces(userID string) ([]domain.Workspace, error) {
	return f.workspaces[userID], nil
}

func (f *fakeWorkspaceService) UpdateWorkspace(string, string, string) (*domain.Workspace, error) {
	panic("not used")
}
func (f *fakeWorkspaceService) DeleteWorkspace(string, string) error { panic("not used") }
func (f *fakeWorkspaceService) AddMember(string, string, string, domain.WorkspaceRole) error {
	panic("not used")
}
func (f *fakeWorkspaceService) RemoveMember(string, string, string) error { panic("not used") }
func (f *fakeWorkspaceService) UpdateMemberRole(string, string, string, domain.WorkspaceRole) error {
	panic("not used")
}
func (f *fakeWorkspaceService) ListMembers(string, string) ([]domain.WorkspaceMember, error) {
	panic("not used")
}
func (f *fakeWorkspaceService) GetMembership(wsID, userID string) (*domain.WorkspaceMember, error) {
	role, ok := f.memberships[wsID+"|"+userID]
	if !ok {
		return nil, nil
	}
	return &domain.WorkspaceMember{WorkspaceID: wsID, UserID: userID, Role: role}, nil
}
func (f *fakeWorkspaceService) WorkspaceExists(workspaceID string) (bool, error) {
	for key := range f.memberships {
		if len(key) > len(workspaceID) && key[:len(workspaceID)+1] == workspaceID+"|" {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeWorkspaceService) CanManage(m *domain.WorkspaceMember) bool {
	return m != nil && (m.Role == domain.WorkspaceOwner || m.Role == domain.WorkspaceAdmin)
}

const (
	wsOne = "11111111-1111-1111-1111-111111111111"
	wsTwo = "22222222-2222-2222-2222-222222222222"
)

// runMiddleware executes the workspace middleware with a pre-seeded auth
// payload and optional X-Workspace-ID header, returning the resolved context
// and the response status (0 when the handler chain was reached).
func runMiddleware(t *testing.T, svc *fakeWorkspaceService, userID, header string) (*WorkspaceContext, int) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	var resolved *WorkspaceContext
	engine.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set(authorizationPayloadKey, &domain.TokenPayload{ID: uuid.New(), UserID: userID, Username: "u"})
		}
		c.Next()
	})
	engine.Use(workspaceMiddleware(svc))
	engine.GET("/x", func(c *gin.Context) {
		resolved = GetWorkspaceContext(c)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	if header != "" {
		req.Header.Set(workspaceHeaderKey, header)
	}
	ctx.Request = req
	engine.ServeHTTP(w, req)

	return resolved, w.Code
}

func TestWorkspaceMiddleware_ExplicitHeader_ResolvesMembership(t *testing.T) {
	svc := newFakeWorkspaceService()
	svc.add("user-1", wsOne, domain.WorkspaceMemberRole)

	resolved, code := runMiddleware(t, svc, "user-1", wsOne)
	assert.Equal(t, http.StatusOK, code)
	assert.NotNil(t, resolved)
	assert.Equal(t, wsOne, resolved.WorkspaceID)
	assert.Equal(t, domain.WorkspaceMemberRole, resolved.Role)
}

// TestWorkspaceMiddleware_NonMember_Forbidden is the core isolation check: a
// caller naming a workspace they do not belong to is denied.
func TestWorkspaceMiddleware_NonMember_Forbidden(t *testing.T) {
	svc := newFakeWorkspaceService()
	svc.add("user-1", wsOne, domain.WorkspaceMemberRole)
	svc.add("user-2", wsTwo, domain.WorkspaceOwner)

	resolved, code := runMiddleware(t, svc, "user-1", wsTwo)
	assert.Equal(t, http.StatusForbidden, code, "user-1 must not adopt workspace two")
	assert.Nil(t, resolved)
}

func TestWorkspaceMiddleware_UnknownWorkspace_Forbidden(t *testing.T) {
	svc := newFakeWorkspaceService()
	svc.add("user-1", wsOne, domain.WorkspaceMemberRole)

	_, code := runMiddleware(t, svc, "user-1", "99999999-9999-9999-9999-999999999999")
	assert.Equal(t, http.StatusForbidden, code)
}

func TestWorkspaceMiddleware_Unauthenticated_Unauthorized(t *testing.T) {
	svc := newFakeWorkspaceService()
	_, code := runMiddleware(t, svc, "", wsOne)
	assert.Equal(t, http.StatusUnauthorized, code)
}

// TestWorkspaceMiddleware_SingleWorkspace_AutoResolves preserves the
// single-workspace home deployment: no header, one membership, resolve it.
func TestWorkspaceMiddleware_SingleWorkspace_AutoResolves(t *testing.T) {
	svc := newFakeWorkspaceService()
	svc.add("user-1", wsOne, domain.WorkspaceOwner)

	resolved, code := runMiddleware(t, svc, "user-1", "")
	assert.Equal(t, http.StatusOK, code)
	assert.NotNil(t, resolved)
	assert.Equal(t, wsOne, resolved.WorkspaceID)
}

// TestWorkspaceMiddleware_NoWorkspace_Forbidden fails closed for a user with
// no workspace membership at all.
func TestWorkspaceMiddleware_NoWorkspace_Forbidden(t *testing.T) {
	svc := newFakeWorkspaceService()
	resolved, code := runMiddleware(t, svc, "user-1", "")
	assert.Equal(t, http.StatusForbidden, code)
	assert.Nil(t, resolved)
}

// TestWorkspaceMiddleware_MultipleWorkspaces_RequiresHeader refuses to guess
// when the caller belongs to more than one workspace.
func TestWorkspaceMiddleware_MultipleWorkspaces_RequiresHeader(t *testing.T) {
	svc := newFakeWorkspaceService()
	svc.add("user-1", wsOne, domain.WorkspaceOwner)
	svc.add("user-1", wsTwo, domain.WorkspaceMemberRole)

	resolved, code := runMiddleware(t, svc, "user-1", "")
	assert.Equal(t, http.StatusBadRequest, code)
	assert.Nil(t, resolved)

	// With an explicit header it resolves to the requested workspace.
	resolved, code = runMiddleware(t, svc, "user-1", wsTwo)
	assert.Equal(t, http.StatusOK, code)
	assert.NotNil(t, resolved)
	assert.Equal(t, wsTwo, resolved.WorkspaceID)
}

func TestRequireWorkspaceRole_EnforcesRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		role     domain.WorkspaceRole
		expected int
	}{
		{domain.WorkspaceOwner, http.StatusOK},
		{domain.WorkspaceAdmin, http.StatusOK},
		{domain.WorkspaceMemberRole, http.StatusForbidden},
		{domain.WorkspaceViewer, http.StatusForbidden},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(w)
		engine.Use(func(c *gin.Context) {
			c.Set(workspaceContextKey, &WorkspaceContext{WorkspaceID: wsOne, Role: tc.role})
			c.Next()
		})
		engine.Use(requireWorkspaceRole(domain.WorkspaceOwner, domain.WorkspaceAdmin))
		engine.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

		engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
		assert.Equal(t, tc.expected, w.Code, "role %s", tc.role)
	}
}

func TestRequireWorkspaceRole_NoContext_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(requireWorkspaceRole(domain.WorkspaceOwner))
	engine.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	assert.Equal(t, http.StatusForbidden, w.Code)
}
