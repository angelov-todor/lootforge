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
	"github.com/angelov-todor/lootforge/core/internal/store"
)

// mockRollStore implements store.RollStore for tests.
type mockRollStore struct {
	sessions     []*models.RollSession
	stats        map[string]int
	createErr    error
	listErr      error
	statsErr     error
	nextID       string
}

func newMockRollStore() *mockRollStore {
	return &mockRollStore{nextID: "roll-1", stats: make(map[string]int)}
}

func (m *mockRollStore) CreateRoll(_ context.Context, groupID string, roll *models.RollSession) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	id := m.nextID
	roll.ID = id
	roll.GroupID = groupID
	m.sessions = append(m.sessions, roll)
	// update in-memory stats
	m.stats[roll.WinnerID]++
	return id, nil
}

func (m *mockRollStore) ListRolls(_ context.Context, groupID string, _ store.ListRollsOpts) ([]*models.RollSession, string, error) {
	if m.listErr != nil {
		return nil, "", m.listErr
	}
	var out []*models.RollSession
	for _, s := range m.sessions {
		if s.GroupID == groupID {
			out = append(out, s)
		}
	}
	return out, "", nil
}

func (m *mockRollStore) GetRollStats(_ context.Context, groupID string) (map[string]int, error) {
	if m.statsErr != nil {
		return nil, m.statsErr
	}
	result := make(map[string]int)
	for _, s := range m.sessions {
		if s.GroupID == groupID {
			result[s.WinnerID]++
		}
	}
	return result, nil
}

func setupRollHandler() (*RollHandler, *mockMemberStore, *mockRollStore, *mockRoleStore, *mockGroupStore) {
	ms := newMockMemberStore()
	rs := newMockRollStore()
	roleS := newMockRoleStore()
	gs := newMockGroupStore()
	h := NewRollHandler(ms, rs, roleS, gs)
	return h, ms, rs, roleS, gs
}

func TestExecuteRoll_OK(t *testing.T) {
	h, ms, _, roleS, gs := setupRollHandler()

	// Set up group with pure_random strategy
	gs.groups["g1"] = &models.Group{
		ID:       "g1",
		Name:     "G1",
		OwnerID:  "user-1",
		Strategy: models.StrategyConfig{Type: models.StrategyPureRandom},
	}

	// Set up members
	ms.members["g1"] = map[string]*models.Member{
		"m1": {ID: "m1", GroupID: "g1", Name: "Alice"},
		"m2": {ID: "m2", GroupID: "g1", Name: "Bob"},
	}

	// Give user view access
	withOwnerRole(roleS, "user-1", "g1")

	body, _ := json.Marshal(RollRequest{
		ParticipantIDs: []string{"m1", "m2"},
		Item:           "Sword of Truth",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/groups/g1/rolls", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.ExecuteRoll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rec.Code, rec.Body.String())
	}

	var got RollResponse
	json.NewDecoder(rec.Body).Decode(&got)
	if got.Winner == nil {
		t.Error("expected a winner")
	}
	if got.Item != "Sword of Truth" {
		t.Errorf("expected item 'Sword of Truth' got %q", got.Item)
	}
	if len(got.Participants) != 2 {
		t.Errorf("expected 2 participants got %d", len(got.Participants))
	}
}

func TestExecuteRoll_EmptyParticipants(t *testing.T) {
	h, _, _, roleS, gs := setupRollHandler()
	gs.groups["g1"] = &models.Group{ID: "g1", Strategy: models.StrategyConfig{Type: models.StrategyPureRandom}}
	withOwnerRole(roleS, "user-1", "g1")

	body, _ := json.Marshal(RollRequest{ParticipantIDs: []string{}, Item: "Item"})
	req := httptest.NewRequest(http.MethodPost, "/api/groups/g1/rolls", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.ExecuteRoll(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}

func TestExecuteRoll_Forbidden(t *testing.T) {
	h, _, _, _, _ := setupRollHandler()

	body, _ := json.Marshal(RollRequest{ParticipantIDs: []string{"m1"}, Item: "Item"})
	req := httptest.NewRequest(http.MethodPost, "/api/groups/g1/rolls", bytes.NewReader(body))
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-99"})

	rec := httptest.NewRecorder()
	h.ExecuteRoll(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestListRolls_OK(t *testing.T) {
	h, _, rs, roleS, _ := setupRollHandler()
	withOwnerRole(roleS, "user-1", "g1")

	// Pre-seed a roll session
	rs.sessions = append(rs.sessions, &models.RollSession{
		ID:       "roll-1",
		GroupID:  "g1",
		WinnerID: "m1",
		Item:     "Sword",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/groups/g1/rolls", nil)
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.ListRolls(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var got listRollsResponse
	json.NewDecoder(rec.Body).Decode(&got)
	if len(got.Rolls) != 1 {
		t.Errorf("expected 1 roll got %d", len(got.Rolls))
	}
}

func TestGetStats_OK(t *testing.T) {
	h, _, rs, roleS, _ := setupRollHandler()
	withOwnerRole(roleS, "user-1", "g1")

	rs.sessions = append(rs.sessions,
		&models.RollSession{ID: "r1", GroupID: "g1", WinnerID: "m1"},
		&models.RollSession{ID: "r2", GroupID: "g1", WinnerID: "m1"},
		&models.RollSession{ID: "r3", GroupID: "g1", WinnerID: "m2"},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/groups/g1/rolls/stats", nil)
	req.SetPathValue("gid", "g1")
	req = authRequest(req, &auth.AuthUser{UID: "user-1"})

	rec := httptest.NewRecorder()
	h.GetStats(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var got statsResponse
	json.NewDecoder(rec.Body).Decode(&got)
	if got.Wins["m1"] != 2 {
		t.Errorf("expected m1 wins=2 got %d", got.Wins["m1"])
	}
	if got.Wins["m2"] != 1 {
		t.Errorf("expected m2 wins=1 got %d", got.Wins["m2"])
	}
}
