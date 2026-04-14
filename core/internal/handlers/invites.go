package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/httputil"
	"github.com/angelov-todor/lootforge/core/internal/middleware"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

type InviteHandler struct {
	inviteStore store.InviteStore
	roleStore   store.RoleStore
	groupStore  store.GroupStore
	memberStore store.MemberStore
}

func NewInviteHandler(
	inviteStore store.InviteStore,
	roleStore store.RoleStore,
	groupStore store.GroupStore,
	memberStore store.MemberStore,
) *InviteHandler {
	return &InviteHandler{
		inviteStore: inviteStore,
		roleStore:   roleStore,
		groupStore:  groupStore,
		memberStore: memberStore,
	}
}

func (h *InviteHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	groupID := r.PathValue("gid")
	if groupID == "" {
		httputil.WriteError(w, 400, "group id is required")
		return
	}

	if err := auth.CanManageGroup(r.Context(), authUser.UID, groupID, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	token, err := generateToken()
	if err != nil {
		httputil.WriteError(w, 500, "failed to generate invite token")
		return
	}

	invite := &models.Invite{
		Token:     token,
		GroupID:   groupID,
		CreatedBy: authUser.UID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := h.inviteStore.CreateInvite(r.Context(), invite); err != nil {
		httputil.WriteError(w, 500, "failed to create invite")
		return
	}

	httputil.WriteJSON(w, 201, invite)
}

func (h *InviteHandler) Accept(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	token := r.PathValue("token")
	if token == "" {
		httputil.WriteError(w, 400, "invite token is required")
		return
	}

	invite, err := h.inviteStore.GetInvite(r.Context(), token)
	if err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "invite not found")
			return
		}
		httputil.WriteError(w, 500, "failed to get invite")
		return
	}

	if time.Now().After(invite.ExpiresAt) {
		_ = h.inviteStore.DeleteInvite(r.Context(), token)
		httputil.WriteError(w, 410, "invite has expired")
		return
	}

	// Check if user already has a role in the group
	if _, err := h.roleStore.GetRole(r.Context(), invite.GroupID, authUser.UID); err == nil {
		httputil.WriteError(w, 409, "already a member of this group")
		return
	}

	role := &models.GroupUserRole{
		UserID:  authUser.UID,
		GroupID: invite.GroupID,
		Role:    models.RoleMember,
	}
	if err := h.roleStore.SetRole(r.Context(), role); err != nil {
		httputil.WriteError(w, 500, "failed to join group")
		return
	}

	// Add as a member entry
	member := &models.Member{
		GroupID:   invite.GroupID,
		Name:      authUser.DisplayName,
		Role:      models.RoleMember,
		CreatedAt: time.Now(),
	}
	if _, err := h.memberStore.AddMember(r.Context(), invite.GroupID, member); err != nil {
		httputil.WriteError(w, 500, "failed to add member")
		return
	}

	group, err := h.groupStore.GetGroup(r.Context(), invite.GroupID)
	if err != nil {
		httputil.WriteError(w, 500, "failed to get group")
		return
	}

	httputil.WriteJSON(w, 200, group)
}

func (h *InviteHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	token := r.PathValue("token")
	if token == "" {
		httputil.WriteError(w, 400, "invite token is required")
		return
	}

	invite, err := h.inviteStore.GetInvite(r.Context(), token)
	if err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "invite not found")
			return
		}
		httputil.WriteError(w, 500, "failed to get invite")
		return
	}

	if err := auth.CanManageGroup(r.Context(), authUser.UID, invite.GroupID, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	if err := h.inviteStore.DeleteInvite(r.Context(), token); err != nil {
		httputil.WriteError(w, 500, "failed to revoke invite")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
