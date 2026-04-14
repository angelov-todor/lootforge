package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/models"
)

// mockInviteStore implements store.InviteStore for tests.
type mockInviteStore struct {
	invites   map[string]*models.Invite
	createErr error
	getErr    error
	deleteErr error
}

func newMockInviteStore() *mockInviteStore {
	return &mockInviteStore{invites: make(map[string]*models.Invite)}
}

func (m *mockInviteStore) CreateInvite(_ context.Context, invite *models.Invite) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.invites[invite.Token] = invite
	return nil
}

func (m *mockInviteStore) GetInvite(_ context.Context, token string) (*models.Invite, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	inv, ok := m.invites[token]
	if !ok {
		return nil, models.ErrNotFound
	}
	return inv, nil
}

func (m *mockInviteStore) DeleteInvite(_ context.Context, token string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.invites, token)
	return nil
}

func newInviteHandler() (*InviteHandler, *mockInviteStore, *mockRoleStore, *mockGroupStore, *mockMemberStore) {
	is := newMockInviteStore()
	rs := newMockRoleStore()
	gs := newMockGroupStore()
	ms := newMockMemberStore()
	h := NewInviteHandler(is, rs, gs, ms)
	return h, is, rs, gs, ms
}

func TestInviteCreate_OK(t *testing.T) {
	h, _, rs, _, _ := newInviteHandler()
	withOwnerRole(rs, "user-1", "group-1")

	req := httptest.NewRequest(http.MethodPost, "/api/groups/group-1/invites", nil)
	req.SetPathValue("gid", "group-1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d: %s", rec.Code, rec.Body.String())
	}

	var got models.Invite
	json.NewDecoder(rec.Body).Decode(&got)
	if got.Token == "" {
		t.Error("expected non-empty token")
	}
	if got.GroupID != "group-1" {
		t.Errorf("expected groupID 'group-1' got %q", got.GroupID)
	}
	if got.ExpiresAt.Before(time.Now()) {
		t.Error("expected future expiry")
	}
}

func TestInviteCreate_Forbidden(t *testing.T) {
	h, _, _, _, _ := newInviteHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/groups/group-1/invites", nil)
	req.SetPathValue("gid", "group-1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"}) // no role

	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestInviteAccept_OK(t *testing.T) {
	h, is, rs, gs, _ := newInviteHandler()

	gs.groups["group-1"] = &models.Group{ID: "group-1", Name: "Test Group", OwnerID: "owner-1"}
	is.invites["test-token"] = &models.Invite{
		Token:     "test-token",
		GroupID:   "group-1",
		CreatedBy: "owner-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	// Owner role exists for the owner, not the joining user
	withOwnerRole(rs, "owner-1", "group-1")

	req := httptest.NewRequest(http.MethodPost, "/api/invites/test-token/accept", nil)
	req.SetPathValue("token", "test-token")
	req = authRequest(req, &auth.AuthUser{UID: "user-2", DisplayName: "User Two"})

	rec := httptest.NewRecorder()
	h.Accept(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rec.Code, rec.Body.String())
	}

	var got models.Group
	json.NewDecoder(rec.Body).Decode(&got)
	if got.ID != "group-1" {
		t.Errorf("expected group-1 got %q", got.ID)
	}

	// Verify role was set
	role, err := rs.GetRole(context.Background(), "group-1", "user-2")
	if err != nil {
		t.Fatalf("role not set: %v", err)
	}
	if role.Role != models.RoleMember {
		t.Errorf("expected member role got %q", role.Role)
	}
}

func TestInviteAccept_Expired(t *testing.T) {
	h, is, _, _, _ := newInviteHandler()

	is.invites["expired-token"] = &models.Invite{
		Token:     "expired-token",
		GroupID:   "group-1",
		CreatedBy: "owner-1",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/invites/expired-token/accept", nil)
	req.SetPathValue("token", "expired-token")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"})

	rec := httptest.NewRecorder()
	h.Accept(rec, req)

	if rec.Code != http.StatusGone {
		t.Fatalf("expected 410 got %d", rec.Code)
	}
}

func TestInviteAccept_NotFound(t *testing.T) {
	h, _, _, _, _ := newInviteHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/invites/nonexistent/accept", nil)
	req.SetPathValue("token", "nonexistent")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"})

	rec := httptest.NewRecorder()
	h.Accept(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rec.Code)
	}
}

func TestInviteAccept_AlreadyMember(t *testing.T) {
	h, is, rs, _, _ := newInviteHandler()

	is.invites["test-token"] = &models.Invite{
		Token:     "test-token",
		GroupID:   "group-1",
		CreatedBy: "owner-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	withMemberRole(rs, "user-2", "group-1")

	req := httptest.NewRequest(http.MethodPost, "/api/invites/test-token/accept", nil)
	req.SetPathValue("token", "test-token")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"})

	rec := httptest.NewRecorder()
	h.Accept(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 got %d", rec.Code)
	}
}

func TestInviteRevoke_OK(t *testing.T) {
	h, is, rs, _, _ := newInviteHandler()

	is.invites["test-token"] = &models.Invite{
		Token:     "test-token",
		GroupID:   "group-1",
		CreatedBy: "owner-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	withOwnerRole(rs, "user-1", "group-1")

	req := httptest.NewRequest(http.MethodDelete, "/api/invites/test-token", nil)
	req.SetPathValue("token", "test-token")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Revoke(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify invite was deleted
	if _, ok := is.invites["test-token"]; ok {
		t.Error("expected invite to be deleted")
	}
}

func TestInviteRevoke_Forbidden(t *testing.T) {
	h, is, _, _, _ := newInviteHandler()

	is.invites["test-token"] = &models.Invite{
		Token:     "test-token",
		GroupID:   "group-1",
		CreatedBy: "owner-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/invites/test-token", nil)
	req.SetPathValue("token", "test-token")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"}) // no role

	rec := httptest.NewRecorder()
	h.Revoke(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}
