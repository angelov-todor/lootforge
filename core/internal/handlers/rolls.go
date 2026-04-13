package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/httputil"
	"github.com/angelov-todor/lootforge/core/internal/middleware"
	"github.com/angelov-todor/lootforge/core/internal/models"
	"github.com/angelov-todor/lootforge/core/internal/store"
	"github.com/angelov-todor/lootforge/core/internal/strategies"
)

type RollHandler struct {
	memberStore store.MemberStore
	rollStore   store.RollStore
	roleStore   store.RoleStore
	groupStore  store.GroupStore
}

func NewRollHandler(
	memberStore store.MemberStore,
	rollStore store.RollStore,
	roleStore store.RoleStore,
	groupStore store.GroupStore,
) *RollHandler {
	return &RollHandler{
		memberStore: memberStore,
		rollStore:   rollStore,
		roleStore:   roleStore,
		groupStore:  groupStore,
	}
}

type RollRequest struct {
	ParticipantIDs []string `json:"participantIDs"`
	Item           string   `json:"item"`
}

type RollResponse struct {
	ID           string               `json:"id"`
	Winner       *models.Member       `json:"winner"`
	Participants []*models.Member     `json:"participants"`
	Item         string               `json:"item"`
	Strategy     models.StrategyConfig `json:"strategy"`
}

func (h *RollHandler) ExecuteRoll(w http.ResponseWriter, r *http.Request) {
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

	// 1. Authorization check
	if err := auth.CanViewGroup(r.Context(), authUser.UID, gid, h.roleStore); err != nil {
		httputil.WriteError(w, 403, "forbidden")
		return
	}

	// 2. Decode request
	var req RollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, 400, "invalid request body")
		return
	}

	// 3. Validate participants
	if len(req.ParticipantIDs) == 0 {
		httputil.WriteError(w, 400, "participantIDs must not be empty")
		return
	}

	// 4. Load group for strategy config
	group, err := h.groupStore.GetGroup(r.Context(), gid)
	if err != nil {
		if err == models.ErrNotFound {
			httputil.WriteError(w, 404, "group not found")
			return
		}
		httputil.WriteError(w, 500, "failed to load group")
		return
	}

	// 5. Load participants from member store
	participants := make([]*models.Member, 0, len(req.ParticipantIDs))
	for _, pid := range req.ParticipantIDs {
		member, err := h.memberStore.GetMember(r.Context(), gid, pid)
		if err != nil {
			if err == models.ErrNotFound {
				httputil.WriteError(w, 404, "participant not found: "+pid)
				return
			}
			httputil.WriteError(w, 500, "failed to load participant")
			return
		}
		participants = append(participants, member)
	}

	// 6. Create strategy
	strategy, err := strategies.NewStrategy(group.Strategy)
	if err != nil {
		httputil.WriteError(w, 400, "unsupported strategy: "+err.Error())
		return
	}

	// 7. Execute roll
	winnerID, err := strategy.Roll(participants)
	if err != nil {
		httputil.WriteError(w, 500, "roll failed: "+err.Error())
		return
	}

	// 8. Separate winner and losers
	var winner *models.Member
	var losers []*models.Member
	for _, p := range participants {
		if p.ID == winnerID {
			winner = p
		} else {
			losers = append(losers, p)
		}
	}
	if winner == nil {
		httputil.WriteError(w, 500, "winner not found among participants")
		return
	}

	// 9. Adjust stats after roll
	if err := strategy.AdjustAfterRoll(winner, losers); err != nil {
		httputil.WriteError(w, 500, "failed to adjust stats after roll")
		return
	}

	// 10. Persist RollSession
	session := &models.RollSession{
		GroupID:        gid,
		Strategy:       group.Strategy,
		ParticipantIDs: req.ParticipantIDs,
		WinnerID:       winnerID,
		Item:           req.Item,
		CreatedAt:      time.Now(),
	}

	sessionID, err := h.rollStore.CreateRoll(r.Context(), gid, session)
	if err != nil {
		httputil.WriteError(w, 500, "failed to save roll session")
		return
	}
	session.ID = sessionID

	// 11. BatchUpdate members with adjusted stats
	allAdjusted := append([]*models.Member{winner}, losers...)
	if err := h.memberStore.BatchUpdateMembers(r.Context(), gid, allAdjusted); err != nil {
		httputil.WriteError(w, 500, "failed to persist member updates")
		return
	}

	// 12. Return response
	httputil.WriteJSON(w, 200, RollResponse{
		ID:           session.ID,
		Winner:       winner,
		Participants: participants,
		Item:         req.Item,
		Strategy:     group.Strategy,
	})
}

type listRollsResponse struct {
	Rolls      []*models.RollSession `json:"rolls"`
	NextCursor string                `json:"nextCursor,omitempty"`
}

func (h *RollHandler) ListRolls(w http.ResponseWriter, r *http.Request) {
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

	q := r.URL.Query()
	opts := store.ListRollsOpts{
		Cursor:   q.Get("cursor"),
		MemberID: q.Get("member"),
		Item:     q.Get("item"),
	}
	if limitStr := q.Get("limit"); limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 {
			opts.Limit = n
		}
	}

	rolls, nextCursor, err := h.rollStore.ListRolls(r.Context(), gid, opts)
	if err != nil {
		httputil.WriteError(w, 500, "failed to list rolls")
		return
	}

	httputil.WriteJSON(w, 200, listRollsResponse{
		Rolls:      rolls,
		NextCursor: nextCursor,
	})
}

type statsResponse struct {
	Wins map[string]int `json:"wins"`
}

func (h *RollHandler) GetStats(w http.ResponseWriter, r *http.Request) {
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

	stats, err := h.rollStore.GetRollStats(r.Context(), gid)
	if err != nil {
		httputil.WriteError(w, 500, "failed to get stats")
		return
	}

	httputil.WriteJSON(w, 200, statsResponse{Wins: stats})
}
