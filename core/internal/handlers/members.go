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

type MemberHandler struct {
	memberStore store.MemberStore
	roleStore   store.RoleStore
}

func NewMemberHandler(memberStore store.MemberStore, roleStore store.RoleStore) *MemberHandler {
	return &MemberHandler{memberStore: memberStore, roleStore: roleStore}
}

func (h *MemberHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	gid := r.PathValue("gid")
	if gid == "" {
		httputil.WriteError(w, 400, "group id is required")
		return
	}

	if err := auth.CanViewGroup(r.Context(), authUser.UID, gid, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	members, err := h.memberStore.ListMembers(r.Context(), gid)
	if err != nil {
		httputil.WriteError(w, 500, "failed to list members")
		return
	}

	httputil.WriteJSON(w, 200, members)
}

type addMemberRequest struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Luck     int    `json:"luck"`
	Priority int    `json:"priority"`
	Points   int    `json:"points"`
}

func (h *MemberHandler) Add(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	gid := r.PathValue("gid")
	if gid == "" {
		httputil.WriteError(w, 400, "group id is required")
		return
	}

	if err := auth.CanManageGroup(r.Context(), authUser.UID, gid, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	var req addMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, 400, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		httputil.WriteError(w, 400, "name is required")
		return
	}

	member := &models.Member{
		GroupID:   gid,
		Name:      req.Name,
		Role:      req.Role,
		Luck:      req.Luck,
		Priority:  req.Priority,
		Points:    req.Points,
		CreatedAt: time.Now(),
	}

	id, err := h.memberStore.AddMember(r.Context(), gid, member)
	if err != nil {
		httputil.WriteError(w, 500, "failed to add member")
		return
	}
	member.ID = id

	httputil.WriteJSON(w, 201, member)
}

func (h *MemberHandler) Update(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	gid := r.PathValue("gid")
	id := r.PathValue("id")
	if gid == "" || id == "" {
		httputil.WriteError(w, 400, "group id and member id are required")
		return
	}

	if err := auth.CanManageGroup(r.Context(), authUser.UID, gid, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	member, err := h.memberStore.GetMember(r.Context(), gid, id)
	if err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "member not found")
			return
		}
		httputil.WriteError(w, 500, "failed to get member")
		return
	}

	var req addMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, 400, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		httputil.WriteError(w, 400, "name is required")
		return
	}

	member.Name = req.Name
	member.Role = req.Role
	member.Luck = req.Luck
	member.Priority = req.Priority
	member.Points = req.Points

	if err := h.memberStore.UpdateMember(r.Context(), gid, member); err != nil {
		httputil.WriteError(w, 500, "failed to update member")
		return
	}

	httputil.WriteJSON(w, 200, member)
}

func (h *MemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	gid := r.PathValue("gid")
	id := r.PathValue("id")
	if gid == "" || id == "" {
		httputil.WriteError(w, 400, "group id and member id are required")
		return
	}

	if err := auth.CanManageGroup(r.Context(), authUser.UID, gid, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	if err := h.memberStore.DeleteMember(r.Context(), gid, id); err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "member not found")
			return
		}
		httputil.WriteError(w, 500, "failed to delete member")
		return
	}

	w.WriteHeader(204)
}

type adjustPointsRequest struct {
	Amount int `json:"amount"`
}

func (h *MemberHandler) AdjustPoints(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.UserFromContext(r.Context())
	if authUser == nil {
		httputil.WriteError(w, 401, "unauthorized")
		return
	}

	gid := r.PathValue("gid")
	id := r.PathValue("id")
	if gid == "" || id == "" {
		httputil.WriteError(w, 400, "group id and member id are required")
		return
	}

	if err := auth.CanManageGroup(r.Context(), authUser.UID, gid, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	member, err := h.memberStore.GetMember(r.Context(), gid, id)
	if err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "member not found")
			return
		}
		httputil.WriteError(w, 500, "failed to get member")
		return
	}

	var req adjustPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, 400, "invalid request body")
		return
	}

	member.Points += req.Amount

	if err := h.memberStore.UpdateMember(r.Context(), gid, member); err != nil {
		httputil.WriteError(w, 500, "failed to update member points")
		return
	}

	httputil.WriteJSON(w, 200, member)
}
