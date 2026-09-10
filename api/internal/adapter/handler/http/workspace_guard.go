package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// resourceInWorkspace guards single-resource handlers against cross-tenant
// access (issue #74 Phase 2).
//
// It is a read of the resource's workspace_id — a value that is immutable
// after creation, so there is no time-of-check/time-of-use window between the
// guard and the operation that follows.
//
// Fails closed: a missing workspace context, an empty resource workspace, or
// a mismatch all return 404 — deliberately the same status as a missing
// resource, so the endpoint is not an existence oracle for other tenants'
// IDs.
func resourceInWorkspace(ctx *gin.Context, resourceWorkspaceID string) bool {
	wc := GetWorkspaceContext(ctx)
	if wc == nil || wc.WorkspaceID == "" || resourceWorkspaceID == "" || resourceWorkspaceID != wc.WorkspaceID {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return false
	}
	return true
}

// filterByWorkspace returns only the rows belonging to the given workspace.
// Used for list endpoints until the Phase 2 query-scoping sweep pushes the
// filter into the repositories.
func filterByWorkspace[T any](rows []T, workspaceID string, workspaceOf func(T) string) []T {
	if workspaceID == "" {
		return nil
	}
	out := make([]T, 0, len(rows))
	for _, r := range rows {
		if workspaceOf(r) == workspaceID {
			out = append(out, r)
		}
	}
	return out
}
