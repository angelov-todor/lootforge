package httputil_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/angelov-todor/lootforge/core/internal/httputil"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"hello": "world"}
	httputil.WriteJSON(w, 200, data)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["hello"] != "world" {
		t.Errorf("expected world, got %s", result["hello"])
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	httputil.WriteError(w, 400, "bad request")

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
	var result httputil.APIError
	json.NewDecoder(w.Body).Decode(&result)
	if result.Message != "bad request" {
		t.Errorf("expected 'bad request', got %s", result.Message)
	}
}
