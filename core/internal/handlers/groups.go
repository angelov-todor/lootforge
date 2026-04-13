package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/httputil"
	"github.com/angelov-todor/lootforge/core/internal/middleware"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

type GroupHandler struct {
	groupStore store.GroupStore
	roleStore  store.RoleStore
}

func NewGroupHandler(groupStore store.GroupStore, roleStore store.RoleStore) *GroupHandler {
	return &GroupHandler{groupStore: groupStore, roleStore: roleStore}
}

type createGroupRequest struct {
	Name     string               `json:"name"`
	Strategy models.StrategyConfig `json:"strategy"`
}

func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, 400, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		httputil.WriteError(w, 400, "name is required")
		return
	}

	group := &models.Group{
		Name:      req.Name,
		OwnerID:   authUser.UID,
		Strategy:  req.Strategy,
		CreatedAt: time.Now(),
	}

	id, err := h.groupStore.CreateGroup(r.Context(), group)
	if err != nil {
		httputil.WriteError(w, 500, "failed to create group")
		return
	}
	group.ID = id

	role := &models.GroupUserRole{
		UserID:  authUser.UID,
		GroupID: id,
		Role:    models.RoleOwner,
	}
	if err := h.roleStore.SetRole(r.Context(), role); err != nil {
		httputil.WriteError(w, 500, "failed to set owner role")
		return
	}

	httputil.WriteJSON(w, 201, group)
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	groups, err := h.groupStore.ListGroupsForUser(r.Context(), authUser.UID)
	if err != nil {
		httputil.WriteError(w, 500, "failed to list groups")
		return
	}

	httputil.WriteJSON(w, 200, groups)
}

func (h *GroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	groupID := r.PathValue("id")
	if groupID == "" {
		httputil.WriteError(w, 400, "group id is required")
		return
	}

	if err := auth.CanViewGroup(r.Context(), authUser.UID, groupID, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	group, err := h.groupStore.GetGroup(r.Context(), groupID)
	if err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "group not found")
			return
		}
		httputil.WriteError(w, 500, "failed to get group")
		return
	}

	httputil.WriteJSON(w, 200, group)
}

func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	groupID := r.PathValue("id")
	if groupID == "" {
		httputil.WriteError(w, 400, "group id is required")
		return
	}

	if err := auth.CanManageGroup(r.Context(), authUser.UID, groupID, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	group, err := h.groupStore.GetGroup(r.Context(), groupID)
	if err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "group not found")
			return
		}
		httputil.WriteError(w, 500, "failed to get group")
		return
	}

	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, 400, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		httputil.WriteError(w, 400, "name is required")
		return
	}

	group.Name = req.Name
	group.Strategy = req.Strategy

	if err := h.groupStore.UpdateGroup(r.Context(), group); err != nil {
		httputil.WriteError(w, 500, "failed to update group")
		return
	}

	httputil.WriteJSON(w, 200, group)
}

func (h *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	groupID := r.PathValue("id")
	if groupID == "" {
		httputil.WriteError(w, 400, "group id is required")
		return
	}

	if err := auth.IsOwner(r.Context(), authUser.UID, groupID, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	if err := h.groupStore.DeleteGroup(r.Context(), groupID); err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "group not found")
			return
		}
		httputil.WriteError(w, 500, "failed to delete group")
		return
	}

	w.WriteHeader(204)
}
