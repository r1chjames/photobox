package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

const (
	// workspaceHeaderKey carries the active workspace ID on authenticated
	// requests (issue #74 §3). Also accepted as a query parameter, since a
	// browser WebSocket upgrade cannot set custom headers.
	workspaceHeaderKey = "X-Workspace-ID"
	workspaceQueryKey  = "workspace_id"
	// workspaceContextKey is the gin context key holding *domain.WorkspaceContext.
	workspaceContextKey = "workspace_context"
)

// WorkspaceContext is the resolved per-request tenancy scope: which workspace
// the caller is acting in, and their role in it. Resolved server-side from
// membership on every request — never carried in the (stateless) token, so a
// removed member loses access immediately (issue #74 §3, D8).
type WorkspaceContext struct {
	WorkspaceID string
	Role        domain.WorkspaceRole
}

// GetWorkspaceContext returns the resolved workspace context, or nil when the
// request was not scoped to a workspace.
func GetWorkspaceContext(ctx *gin.Context) *WorkspaceContext {
	v, exists := ctx.Get(workspaceContextKey)
	if !exists {
		return nil
	}
	wc, _ := v.(*WorkspaceContext)
	return wc
}

// workspaceMiddleware resolves the caller's workspace and enforces
// membership. It must run after authMiddleware.
//
// Resolution order:
//  1. X-Workspace-ID header (or workspace_id query param for WS upgrades):
//     the caller's membership is looked up; no membership → 403.
//  2. Absent header: if the caller belongs to exactly one workspace it is
//     selected automatically. This preserves the single-workspace home
//     deployment and existing clients without weakening isolation — the
//     workspace is still chosen only from the caller's own memberships.
//  3. Absent header with multiple memberships: ambiguous → 400, the client
//     must choose explicitly.
//
// Fail-closed: no workspace context on the request is always an error; a
// caller can never act outside a workspace they belong to.
func workspaceMiddleware(wsSvc port.WorkspaceService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := GetAuthPayload(ctx)
		if payload == nil || payload.UserID == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		requested := ctx.GetHeader(workspaceHeaderKey)
		if requested == "" {
			requested = ctx.Query(workspaceQueryKey)
		}

		if requested != "" {
			membership, err := wsSvc.GetMembership(requested, payload.UserID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve workspace"})
				return
			}
			if membership == nil {
				// Not a member (or the workspace does not exist): deny. 403 for
				// a known-but-unauthorized workspace, deliberately without
				// revealing whether it exists.
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No access to workspace"})
				return
			}
			ctx.Set(workspaceContextKey, &WorkspaceContext{WorkspaceID: requested, Role: membership.Role})
			ctx.Next()
			return
		}

		// No explicit workspace: resolve from the caller's memberships.
		workspaces, err := wsSvc.ListWorkspaces(payload.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve workspace"})
			return
		}
		switch len(workspaces) {
		case 0:
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No workspace available for user"})
			return
		case 1:
			membership, err := wsSvc.GetMembership(workspaces[0].ID, payload.UserID)
			if err != nil || membership == nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No access to workspace"})
				return
			}
			ctx.Set(workspaceContextKey, &WorkspaceContext{WorkspaceID: workspaces[0].ID, Role: membership.Role})
			ctx.Next()
			return
		default:
			// Ambiguous: the client must send X-Workspace-ID.
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "X-Workspace-ID header required when the user belongs to multiple workspaces"})
			return
		}
	}
}

// requireWorkspaceRole enforces the caller's workspace role. Must run after
// workspaceMiddleware. A workspace role is orthogonal to the deployment-level
// user role checked by requireRole.
func requireWorkspaceRole(roles ...domain.WorkspaceRole) gin.HandlerFunc {
	allowed := make(map[domain.WorkspaceRole]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(ctx *gin.Context) {
		wc := GetWorkspaceContext(ctx)
		if wc == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No workspace context"})
			return
		}
		if !allowed[wc.Role] {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient workspace permissions"})
			return
		}
		ctx.Next()
	}
}
