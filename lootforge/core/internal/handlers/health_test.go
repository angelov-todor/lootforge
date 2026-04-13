package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/handlers"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	handlers.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	expected := `{"status":"ok"}`
	if body := w.Body.String(); body != expected+"\n" {
		t.Errorf("expected body %q, got %q", expected, body)
	}
}
