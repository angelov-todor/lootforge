package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/angelov-todor/lootforge/core/internal/httputil"
	"github.com/angelov-todor/lootforge/core/internal/middleware"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

type UserHandler struct {
	userStore store.UserStore
}

func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	now := time.Now()
	user := &models.User{
		ID:          authUser.UID,
		Email:       authUser.Email,
		DisplayName: authUser.DisplayName,
		PhotoURL:    authUser.PhotoURL,
		CreatedAt:   now,
		LastLoginAt: now,
	}

	if err := h.userStore.UpsertUser(r.Context(), user); err != nil {
		log.Printf("ERROR UpsertUser uid=%s: %v", authUser.UID, err)
		httputil.WriteError(w, 500, "failed to sync user")
		return
	}

	httputil.WriteJSON(w, 200, user)
}
