package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

// WorkspaceService implements workspace creation, membership management, and
// role checks (issue #74). Ownership is data (workspace_members), never a
// token claim.
type WorkspaceService struct {
	repo port.WorkspaceRepository
}

func NewWorkspaceService(repo port.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{repo: repo}
}

var slugRe = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$`)

// CreateWorkspace creates a workspace owned by userID with the creator as an
// owner member. The slug defaults to a slugified name when empty.
func (s *WorkspaceService) CreateWorkspace(name, slug, userID string) (*domain.Workspace, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("workspace name is required")
	}
	if slug == "" {
		slug = slugify(name)
	}
	if !slugRe.MatchString(slug) {
		return nil, fmt.Errorf("invalid workspace slug %q (lowercase letters, digits, hyphens)", slug)
	}

	w := &domain.Workspace{
		ID:                uuid.New().String(),
		Name:              name,
		Slug:              slug,
		OwnerUserID:       userID,
		StorageUsedBytes:  0,
		StorageLimitBytes: 0, // 0 = unlimited until quotas land
	}
	err := s.repo.CreateWorkspace(w)
	if err != nil {
		// Surface duplicate-slug as a friendly conflict.
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, fmt.Errorf("workspace slug %q is already taken", slug)
		}
		return nil, err
	}
	if err := s.repo.AddMember(w.ID, userID, domain.WorkspaceOwner); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WorkspaceService) GetWorkspace(workspaceID, userID string) (*domain.Workspace, error) {
	m, err := s.repo.GetMembership(workspaceID, userID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, domain.ErrForbidden
	}
	return s.repo.GetWorkspaceById(workspaceID)
}

func (s *WorkspaceService) ListWorkspaces(userID string) ([]domain.Workspace, error) {
	return s.repo.ListWorkspacesForUser(userID)
}

func (s *WorkspaceService) UpdateWorkspace(workspaceID, userID, name string) (*domain.Workspace, error) {
	m, err := s.repo.GetMembership(workspaceID, userID)
	if err != nil {
		return nil, err
	}
	if m == nil || !s.CanManage(m) {
		return nil, domain.ErrForbidden
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("workspace name is required")
	}
	w, err := s.repo.GetWorkspaceById(workspaceID)
	if err != nil {
		return nil, err
	}
	w.Name = strings.TrimSpace(name)
	err = s.repo.UpdateWorkspace(w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WorkspaceService) DeleteWorkspace(workspaceID, userID string) error {
	m, err := s.repo.GetMembership(workspaceID, userID)
	if err != nil {
		return err
	}
	if m == nil || m.Role != domain.WorkspaceOwner {
		return domain.ErrForbidden
	}
	return s.repo.DeleteWorkspace(workspaceID)
}

func (s *WorkspaceService) AddMember(workspaceID, actorUserID, targetUserID string, role domain.WorkspaceRole) error {
	actor, err := s.repo.GetMembership(workspaceID, actorUserID)
	if err != nil {
		return err
	}
	if actor == nil || !s.CanManage(actor) {
		return domain.ErrForbidden
	}
	return s.repo.AddMember(workspaceID, targetUserID, role)
}

func (s *WorkspaceService) RemoveMember(workspaceID, actorUserID, targetUserID string) error {
	actor, err := s.repo.GetMembership(workspaceID, actorUserID)
	if err != nil {
		return err
	}
	if actor == nil || !s.CanManage(actor) {
		return domain.ErrForbidden
	}

	// An owner removing another owner is only allowed when a different owner
	// remains (last-owner rule, issue #74 open question).
	target, err := s.repo.GetMembership(workspaceID, targetUserID)
	if err != nil {
		return err
	}
	if target == nil {
		return domain.ErrDataNotFound
	}
	// Last-owner rule (issue #74 open question): an owner may not be removed
	// while no other owner remains — including when they remove themselves,
	// which would otherwise orphan the workspace.
	if target.Role == domain.WorkspaceOwner {
		remaining, err := s.countOwners(workspaceID, targetUserID)
		if err != nil {
			return err
		}
		if remaining == 0 {
			return errors.New("cannot remove the last workspace owner")
		}
	}
	return s.repo.RemoveMember(workspaceID, targetUserID)
}

func (s *WorkspaceService) UpdateMemberRole(workspaceID, actorUserID, targetUserID string, role domain.WorkspaceRole) error {
	actor, err := s.repo.GetMembership(workspaceID, actorUserID)
	if err != nil {
		return err
	}
	if actor == nil || !s.CanManage(actor) {
		return domain.ErrForbidden
	}
	return s.repo.UpdateMemberRole(workspaceID, targetUserID, role)
}

func (s *WorkspaceService) ListMembers(workspaceID, userID string) ([]domain.WorkspaceMember, error) {
	m, err := s.repo.GetMembership(workspaceID, userID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, domain.ErrForbidden
	}
	return s.repo.ListMembers(workspaceID)
}

func (s *WorkspaceService) GetMembership(workspaceID, userID string) (*domain.WorkspaceMember, error) {
	return s.repo.GetMembership(workspaceID, userID)
}

func (s *WorkspaceService) CanManage(m *domain.WorkspaceMember) bool {
	if m == nil {
		return false
	}
	return m.Role == domain.WorkspaceOwner || m.Role == domain.WorkspaceAdmin
}

// countOwners returns the number of owners of workspaceID excluding userID.
func (s *WorkspaceService) countOwners(workspaceID, excludeUserID string) (int, error) {
	members, err := s.repo.ListMembers(workspaceID)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, m := range members {
		if m.Role == domain.WorkspaceOwner && m.UserID != excludeUserID {
			n++
		}
	}
	return n, nil
}

// slugify converts a display name into a URL-safe slug.
func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		case r == ' ', r == '-', r == '_':
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		default:
			// drop non-ASCII characters
			prevDash = false
		}
	}
	out := strings.TrimSuffix(b.String(), "-")
	if out == "" {
		out = "workspace"
	}
	return out
}
