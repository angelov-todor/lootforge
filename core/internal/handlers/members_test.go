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

// mockMemberStore implements store.MemberStore for tests.
type mockMemberStore struct {
	members         map[string]map[string]*models.Member // groupID -> memberID -> member
	addErr          error
	getErr          error
	listErr         error
	updateErr       error
	deleteErr       error
	batchUpdateErr  error
	nextID          string
}

func newMockMemberStore() *mockMemberStore {
	return &mockMemberStore{
		members: make(map[string]map[string]*models.Member),
		nextID:  "member-1",
	}
}

func (m *mockMemberStore) ListMembers(_ context.Context, groupID string) ([]*models.Member, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var out []*models.Member
	for _, mem := range m.members[groupID] {
		out = append(out, mem)
	}
	return out, nil
}

func (m *mockMemberStore) AddMember(_ context.Context, groupID string, member *models.Member) (string, error) {
	if m.addErr != nil {
		return "", m.addErr
	}
	if m.members[groupID] == nil {
		m.members[groupID] = make(map[string]*models.Member)
	}
	id := m.nextID
	member.ID = id
	m.members[groupID][id] = member
	return id, nil
}

func (m *mockMemberStore) GetMember(_ context.Context, groupID, memberID string) (*models.Member, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	mem, ok := m.members[groupID][memberID]
	if !ok {
		return nil, models.ErrNotFound
	}
	return mem, nil
}

func (m *mockMemberStore) UpdateMember(_ context.Context, groupID string, member *models.Member) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if m.members[groupID] == nil {
		m.members[groupID] = make(map[string]*models.Member)
	}
	m.members[groupID][member.ID] = member
	return nil
}

func (m *mockMemberStore) DeleteMember(_ context.Context, groupID, memberID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.members[groupID][memberID]; !ok {
		return models.ErrNotFound
	}
	delete(m.members[groupID], memberID)
	return nil
}

func (m *mockMemberStore) BatchUpdateMembers(_ context.Context, groupID string, members []*models.Member) error {
	if m.batchUpdateErr != nil {
		return m.batchUpdateErr
	}
	if m.members[groupID] == nil {
		m.members[groupID] = make(map[string]*models.Member)
	}
	for _, mem := range members {
		m.members[groupID][mem.ID] = mem
	}
	return nil
}

func TestMemberList_OK(t *testing.T) {
	ms := newMockMemberStore()
	ms.members["g1"] = map[string]*models.Member{
		"m1": {ID: "m1", GroupID: "g1", Name: "Alice"},
	}
	rs := newMockRoleStore()
	withOwnerRole(rs, "user-1", "g1")
	h := NewMemberHandler(ms, rs)

	req := httptest.NewRequest(http.MethodGet, "/api/groups/g1/members", nil)
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var got []*models.Member
	json.NewDecoder(rec.Body).Decode(&got)
	if len(got) != 1 {
		t.Errorf("expected 1 member got %d", len(got))
	}
}

func TestMemberList_Forbidden(t *testing.T) {
	ms := newMockMemberStore()
	rs := newMockRoleStore()
	h := NewMemberHandler(ms, rs)

	req := httptest.NewRequest(http.MethodGet, "/api/groups/g1/members", nil)
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-99"})

	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestMemberAdd_OK(t *testing.T) {
	ms := newMockMemberStore()
	rs := newMockRoleStore()
	withOwnerRole(rs, "user-1", "g1")
	h := NewMemberHandler(ms, rs)

	body, _ := json.Marshal(map[string]any{"name": "Bob"})
	req := httptest.NewRequest(http.MethodPost, "/api/groups/g1/members", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Add(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d: %s", rec.Code, rec.Body.String())
	}
	var got models.Member
	json.NewDecoder(rec.Body).Decode(&got)
	if got.Name != "Bob" {
		t.Errorf("expected name Bob got %q", got.Name)
	}
}

func TestMemberAdd_Forbidden(t *testing.T) {
	ms := newMockMemberStore()
	rs := newMockRoleStore()
	withMemberRole(rs, "user-2", "g1")
	h := NewMemberHandler(ms, rs)

	body, _ := json.Marshal(map[string]any{"name": "Bob"})
	req := httptest.NewRequest(http.MethodPost, "/api/groups/g1/members", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"})

	rec := httptest.NewRecorder()
	h.Add(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestMemberUpdate_OK(t *testing.T) {
	ms := newMockMemberStore()
	ms.members["g1"] = map[string]*models.Member{
		"m1": {ID: "m1", GroupID: "g1", Name: "Alice", Points: 10},
	}
	rs := newMockRoleStore()
	withOwnerRole(rs, "user-1", "g1")
	h := NewMemberHandler(ms, rs)

	body, _ := json.Marshal(map[string]any{"name": "Alicia", "points": 20})
	req := httptest.NewRequest(http.MethodPut, "/api/groups/g1/members/m1", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req.SetPathValue("id", "m1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rec.Code, rec.Body.String())
	}
	var got models.Member
	json.NewDecoder(rec.Body).Decode(&got)
	if got.Name != "Alicia" {
		t.Errorf("expected name Alicia got %q", got.Name)
	}
}

func TestMemberDelete_OK(t *testing.T) {
	ms := newMockMemberStore()
	ms.members["g1"] = map[string]*models.Member{
		"m1": {ID: "m1", GroupID: "g1", Name: "Alice"},
	}
	rs := newMockRoleStore()
	withOwnerRole(rs, "user-1", "g1")
	h := NewMemberHandler(ms, rs)

	req := httptest.NewRequest(http.MethodDelete, "/api/groups/g1/members/m1", nil)
	req.SetPathValue("gid", "g1")
	req.SetPathValue("id", "m1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", rec.Code)
	}
}

func TestMemberAdjustPoints_OK(t *testing.T) {
	ms := newMockMemberStore()
	ms.members["g1"] = map[string]*models.Member{
		"m1": {ID: "m1", GroupID: "g1", Name: "Alice", Points: 10},
	}
	rs := newMockRoleStore()
	withOwnerRole(rs, "user-1", "g1")
	h := NewMemberHandler(ms, rs)

	body, _ := json.Marshal(map[string]any{"amount": 5})
	req := httptest.NewRequest(http.MethodPatch, "/api/groups/g1/members/m1/points", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req.SetPathValue("id", "m1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.AdjustPoints(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rec.Code, rec.Body.String())
	}
	var got models.Member
	json.NewDecoder(rec.Body).Decode(&got)
	if got.Points != 15 {
		t.Errorf("expected points 15 got %d", got.Points)
	}
}

func TestMemberAdjustPoints_Forbidden(t *testing.T) {
	ms := newMockMemberStore()
	ms.members["g1"] = map[string]*models.Member{
		"m1": {ID: "m1", GroupID: "g1", Name: "Alice", Points: 10},
	}
	rs := newMockRoleStore()
	withMemberRole(rs, "user-2", "g1")
	h := NewMemberHandler(ms, rs)

	body, _ := json.Marshal(map[string]any{"amount": 5})
	req := httptest.NewRequest(http.MethodPatch, "/api/groups/g1/members/m1/points", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req.SetPathValue("id", "m1")
	req = authRequest(req, &auth.AuthUser{UID: "user-2"})

	rec := httptest.NewRecorder()
	h.AdjustPoints(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}
