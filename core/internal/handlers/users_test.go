package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/middleware"
	"github.com/angelov-todor/lootforge/core/internal/models"
)

// mockUserStore implements store.UserStore for tests.
type mockUserStore struct {
	user    *models.User
	upsertErr error
	getErr  error
}

func (m *mockUserStore) UpsertUser(_ context.Context, user *models.User) error {
	if m.upsertErr != nil {
		return m.upsertErr
	}
	m.user = user
	return nil
}

func (m *mockUserStore) GetUser(_ context.Context, id string) (*models.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.user != nil && m.user.ID == id {
		return m.user, nil
	}
	return nil, models.ErrNotFound
}

func addAuthUser(r *http.Request, u *auth.AuthUser) *http.Request {
	verifier := &auth.MockTokenVerifier{User: u}
	var captured *http.Request
	handler := middleware.Auth(verifier)(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		captured = req
	}))
	rec := httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer test-token")
	handler.ServeHTTP(rec, r)
	if captured != nil {
		return captured
	}
	return r
}

func TestGetMe_OK(t *testing.T) {
	store := &mockUserStore{}
	h := NewUserHandler(store)

	authUser := &auth.AuthUser{
		UID:         "user-1",
		Email:       "test@example.com",
		DisplayName: "Test User",
		PhotoURL:    "https://example.com/photo.jpg",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req = addAuthUser(req, authUser)

	rec := httptest.NewRecorder()
	h.GetMe(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}

	var got models.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.ID != authUser.UID {
		t.Errorf("expected ID %q got %q", authUser.UID, got.ID)
	}
	if got.Email != authUser.Email {
		t.Errorf("expected email %q got %q", authUser.Email, got.Email)
	}
	if got.DisplayName != authUser.DisplayName {
		t.Errorf("expected displayName %q got %q", authUser.DisplayName, got.DisplayName)
	}
}

func TestGetMe_Unauthorized(t *testing.T) {
	store := &mockUserStore{}
	h := NewUserHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	rec := httptest.NewRecorder()
	h.GetMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
}

func TestGetMe_UpsertError(t *testing.T) {
	store := &mockUserStore{upsertErr: models.ErrConflict}
	h := NewUserHandler(store)

	authUser := &auth.AuthUser{UID: "user-1", Email: "test@example.com"}
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req = addAuthUser(req, authUser)

	rec := httptest.NewRecorder()
	h.GetMe(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rec.Code)
	}
}
