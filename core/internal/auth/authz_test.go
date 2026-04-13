package auth_test

import (
	"context"
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

type mockRoleStore struct {
	role *models.GroupUserRole
	err  error
}

func (m *mockRoleStore) SetRole(ctx context.Context, role *models.GroupUserRole) error { return nil }
func (m *mockRoleStore) GetRole(ctx context.Context, groupID, userID string) (*models.GroupUserRole, error) {
	return m.role, m.err
}
func (m *mockRoleStore) ListRolesForGroup(ctx context.Context, groupID string) ([]*models.GroupUserRole, error) {
	return nil, nil
}

// Verify mockRoleStore implements store.RoleStore
var _ store.RoleStore = (*mockRoleStore)(nil)

func TestCanManageGroup_Owner(t *testing.T) {
	rs := &mockRoleStore{role: &models.GroupUserRole{Role: models.RoleOwner}}
	err := auth.CanManageGroup(context.Background(), "u1", "g1", rs)
	if err != nil {
		t.Errorf("owner should be able to manage: %v", err)
	}
}

func TestCanManageGroup_Admin(t *testing.T) {
	rs := &mockRoleStore{role: &models.GroupUserRole{Role: models.RoleAdmin}}
	err := auth.CanManageGroup(context.Background(), "u1", "g1", rs)
	if err != nil {
		t.Errorf("admin should be able to manage: %v", err)
	}
}

func TestCanManageGroup_Member(t *testing.T) {
	rs := &mockRoleStore{role: &models.GroupUserRole{Role: models.RoleMember}}
	err := auth.CanManageGroup(context.Background(), "u1", "g1", rs)
	if err == nil {
		t.Error("member should not be able to manage")
	}
}

func TestCanViewGroup_Member(t *testing.T) {
	rs := &mockRoleStore{role: &models.GroupUserRole{Role: models.RoleMember}}
	err := auth.CanViewGroup(context.Background(), "u1", "g1", rs)
	if err != nil {
		t.Errorf("member should be able to view: %v", err)
	}
}

func TestCanViewGroup_NonMember(t *testing.T) {
	rs := &mockRoleStore{err: models.ErrNotFound}
	err := auth.CanViewGroup(context.Background(), "u1", "g1", rs)
	if err == nil {
		t.Error("non-member should not be able to view")
	}
}
