package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

// WorkspaceHandler exposes workspace management (issue #74). All routes are
// authenticated; membership and role checks are enforced in the service.
type WorkspaceHandler struct {
	svc port.WorkspaceService
}

// NewWorkspaceHandler creates a WorkspaceHandler.
func NewWorkspaceHandler(svc port.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{svc: svc}
}

type createWorkspaceRequest struct {
	Name string `json:"name" binding:"required" example:"Holiday Photos"`
	Slug string `json:"slug" example:"holiday-photos"`
}

// Create godoc
//
//	@Summary	Create a workspace
//	@Description	Creates a workspace owned by the caller, who becomes its owner.
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) Create(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var req createWorkspaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	ws, err := wh.svc.CreateWorkspace(req.Name, req.Slug, payload.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, newWorkspaceResponse(ws))
}

// List godoc
//
//	@Summary	List the caller's workspaces
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) List(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	workspaces, err := wh.svc.ListWorkspaces(payload.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	resp := make([]workspaceResponse, 0, len(workspaces))
	for i := range workspaces {
		resp = append(resp, newWorkspaceResponse(&workspaces[i]))
	}
	handleSuccess(ctx, resp)
}

type workspaceIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

// Get godoc
//
//	@Summary	Get a workspace
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) Get(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var req workspaceIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	ws, err := wh.svc.GetWorkspace(req.ID, payload.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, newWorkspaceResponse(ws))
}

type updateWorkspaceRequest struct {
	Name string `json:"name" binding:"required"`
}

// Update godoc
//
//	@Summary	Rename a workspace
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) Update(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var uri workspaceIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		validationError(ctx, err)
		return
	}
	var req updateWorkspaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	ws, err := wh.svc.UpdateWorkspace(uri.ID, payload.UserID, req.Name)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, newWorkspaceResponse(ws))
}

// Delete godoc
//
//	@Summary	Delete a workspace (owner only)
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) Delete(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var req workspaceIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := wh.svc.DeleteWorkspace(req.ID, payload.UserID); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, gin.H{"message": "Workspace deleted"})
}

// ListMembers godoc
//
//	@Summary	List workspace members
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) ListMembers(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var req workspaceIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	members, err := wh.svc.ListMembers(req.ID, payload.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, members)
}

type addMemberRequest struct {
	UserID string               `json:"userId" binding:"required"`
	Role   domain.WorkspaceRole `json:"role" binding:"required,oneof=owner admin member viewer"`
}

// AddMember godoc
//
//	@Summary	Add a member to a workspace
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) AddMember(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var uri workspaceIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		validationError(ctx, err)
		return
	}
	var req addMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := wh.svc.AddMember(uri.ID, payload.UserID, req.UserID, req.Role); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, gin.H{"message": "Member added"})
}

type memberRequest struct {
	UserID string `uri:"userId" binding:"required"`
}

// RemoveMember godoc
//
//	@Summary	Remove a member from a workspace
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) RemoveMember(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var uri workspaceIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		validationError(ctx, err)
		return
	}
	var member memberRequest
	if err := ctx.ShouldBindUri(&member); err != nil {
		validationError(ctx, err)
		return
	}
	if err := wh.svc.RemoveMember(uri.ID, payload.UserID, member.UserID); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, gin.H{"message": "Member removed"})
}

type updateMemberRoleRequest struct {
	Role domain.WorkspaceRole `json:"role" binding:"required,oneof=owner admin member viewer"`
}

// UpdateMemberRole godoc
//
//	@Summary	Change a member's workspace role
//	@Tags		Workspaces
//	@Security	BearerAuth
func (wh *WorkspaceHandler) UpdateMemberRole(ctx *gin.Context) {
	payload := GetAuthPayload(ctx)
	if payload == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var uri workspaceIDRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		validationError(ctx, err)
		return
	}
	var member memberRequest
	if err := ctx.ShouldBindUri(&member); err != nil {
		validationError(ctx, err)
		return
	}
	var req updateMemberRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := wh.svc.UpdateMemberRole(uri.ID, payload.UserID, member.UserID, req.Role); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, gin.H{"message": "Member role updated"})
}
