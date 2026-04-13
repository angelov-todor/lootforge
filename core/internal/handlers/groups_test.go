package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/models"
)

// mockGroupStore implements store.GroupStore for tests.
type mockGroupStore struct {
	groups    map[string]*models.Group
	createErr error
	getErr    error
	listErr   error
	updateErr error
	deleteErr error
	nextID    string
}

func newMockGroupStore() *mockGroupStore {
	return &mockGroupStore{
		groups: make(map[string]*models.Group),
		nextID: "group-1",
	}
}

func (m *mockGroupStore) CreateGroup(_ context.Context, group *models.Group) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	id := m.nextID
	group.ID = id
	m.groups[id] = group
	return id, nil
}

func (m *mockGroupStore) GetGroup(_ context.Context, id string) (*models.Group, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	g, ok := m.groups[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return g, nil
}

func (m *mockGroupStore) ListGroupsForUser(_ context.Context, _ string) ([]*models.Group, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var out []*models.Group
	for _, g := range m.groups {
		out = append(out, g)
	}
	return out, nil
}

func (m *mockGroupStore) UpdateGroup(_ context.Context, group *models.Group) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupStore) DeleteGroup(_ context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.groups[id]; !ok {
		return models.ErrNotFound
	}
	delete(m.groups, id)
	return nil
}

// mockRoleStore implements store.RoleStore for tests.
type mockRoleStore struct {
	roles     map[string]*models.GroupUserRole // key: groupID+":"+userID
	setErr    error
	getErr    error
	listErr   error
}

func newMockRoleStore() *mockRoleStore {
	return &mockRoleStore{roles: make(map[string]*models.GroupUserRole)}
}

func roleKey(groupID, userID string) string { return groupID + ":" + userID }

func (m *mockRoleStore) SetRole(_ context.Context, role *models.GroupUserRole) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.roles[roleKey(role.GroupID, role.UserID)] = role
	return nil
}

func (m *mockRoleStore) GetRole(_ context.Context, groupID, userID string) (*models.GroupUserRole, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	r, ok := m.roles[roleKey(groupID, userID)]
	if !ok {
		return nil, models.ErrNotFound
	}
	return r, nil
}

func (m *mockRoleStore) ListRolesForGroup(_ context.Context, groupID string) ([]*models.GroupUserRole, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var out []*models.GroupUserRole
	for _, r := range m.roles {
		if r.GroupID == groupID {
			out = append(out, r)
		}
	}
	return out, nil
}

// withOwnerRole sets a user as owner in the role store for a given group.
func withOwnerRole(rs *mockRoleStore, userID, groupID string) {
	rs.roles[roleKey(groupID, userID)] = &models.GroupUserRole{
		UserID:  userID,
		GroupID: groupID,
		Role:    models.RoleOwner,
	}
}

// withMemberRole sets a user as member.
func withMemberRole(rs *mockRoleStore, userID, groupID string) {
	rs.roles[roleKey(groupID, userID)] = &models.GroupUserRole{
		UserID:  userID,
		GroupID: groupID,
		Role:    models.RoleMember,
	}
}

func authRequest(r *http.Request, u *auth.AuthUser) *http.Request {
	return addAuthUser(r, u)
}

func TestGroupCreate_OK(t *testing.T) {
	gs := newMockGroupStore()
	rs := newMockRoleStore()
	h := NewGroupHandler(gs, rs)

	body, _ := json.Marshal(map[string]any{"name": "My Group", "strategy": map[string]any{"type": "pure_random"}})
	req := httptest.NewRequest(http.MethodPost, "/api/groups", bytes.NewReader(body))
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d: %s", rec.Code, rec.Body.String())
	}

	var got models.Group
	json.NewDecoder(rec.Body).Decode(&got)
	if got.Name != "My Group" {
		t.Errorf("expected name 'My Group' got %q", got.Name)
	}
	if got.OwnerID != "user-1" {
		t.Errorf("expected ownerID 'user-1' got %q", got.OwnerID)
	}
	// Verify owner role was set
	role, err := rs.GetRole(context.Background(), got.ID, "user-1")
	if err != nil {
		t.Fatalf("role not set: %v", err)
	}
	if role.Role != models.RoleOwner {
		t.Errorf("expected role owner got %q", role.Role)
	}
}

func TestGroupCreate_EmptyName(t *testing.T) {
	gs := newMockGroupStore()
	rs := newMockRoleStore()
	h := NewGroupHandler(gs, rs)

	body, _ := json.Marshal(map[string]any{"name": "  "})
	req := httptest.NewRequest(http.MethodPost, "/api/groups", bytes.NewReader(body))
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}

func TestGroupList_OK(t *testing.T) {
	gs := newMockGroupStore()
	gs.groups["g1"] = &models.Group{ID: "g1", Name: "G1", OwnerID: "user-1"}
	rs := newMockRoleStore()
	h := NewGroupHandler(gs, rs)

	req := httptest.NewRequest(http.MethodGet, "/api/groups", nil)
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var got []*models.Group
	json.NewDecoder(rec.Body).Decode(&got)
	if len(got) != 1 {
		t.Errorf("expected 1 group got %d", len(got))
	}
}

func TestGroupGet_Forbidden(t *testing.T) {
	gs := newMockGroupStore()
	gs.groups["g1"] = &models.Group{ID: "g1", Name: "G1", OwnerID: "owner-1"}
	rs := newMockRoleStore()
	h := NewGroupHandler(gs, rs)

	req := httptest.NewRequest(http.MethodGet, "/api/groups/g1", nil)
	req.SetPathValue("id", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"}) // no role

	rec := httptest.NewRecorder()
	h.Get(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestGroupGet_OK(t *testing.T) {
	gs := newMockGroupStore()
	gs.groups["g1"] = &models.Group{ID: "g1", Name: "G1", OwnerID: "user-1"}
	rs := newMockRoleStore()
	withOwnerRole(rs, "user-1", "g1")
	h := NewGroupHandler(gs, rs)

	req := httptest.NewRequest(http.MethodGet, "/api/groups/g1", nil)
	req.SetPathValue("id", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Get(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestGroupUpdate_Forbidden_Member(t *testing.T) {
	gs := newMockGroupStore()
	gs.groups["g1"] = &models.Group{ID: "g1", Name: "G1", OwnerID: "owner-1"}
	rs := newMockRoleStore()
	withMemberRole(rs, "user-2", "g1")
	h := NewGroupHandler(gs, rs)

	body, _ := json.Marshal(map[string]any{"name": "New Name"})
	req := httptest.NewRequest(http.MethodPut, "/api/groups/g1", bytes.NewReader(body))
	req.SetPathValue("id", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"})

	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestGroupDelete_NotOwner(t *testing.T) {
	gs := newMockGroupStore()
	gs.groups["g1"] = &models.Group{ID: "g1", Name: "G1", OwnerID: "owner-1"}
	rs := newMockRoleStore()
	withMemberRole(rs, "user-2", "g1")
	h := NewGroupHandler(gs, rs)

	req := httptest.NewRequest(http.MethodDelete, "/api/groups/g1", nil)
	req.SetPathValue("id", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"})

	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestGroupDelete_OK(t *testing.T) {
	gs := newMockGroupStore()
	gs.groups["g1"] = &models.Group{ID: "g1", Name: "G1", OwnerID: "user-1"}
	rs := newMockRoleStore()
	withOwnerRole(rs, "user-1", "g1")
	h := NewGroupHandler(gs, rs)

	req := httptest.NewRequest(http.MethodDelete, "/api/groups/g1", nil)
	req.SetPathValue("id", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", rec.Code)
	}
}
