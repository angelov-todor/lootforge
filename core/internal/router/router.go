package router

import (
	"net/http"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/handlers"
	"github.com/angelov-todor/lootforge/core/internal/middleware"
	"github.com/angelov-todor/lootforge/core/internal/store"
)

type Deps struct {
	Verifier    auth.TokenVerifier
	UserStore   store.UserStore
	GroupStore  store.GroupStore
	RoleStore   store.RoleStore
	MemberStore store.MemberStore
	RollStore   store.RollStore
	InviteStore store.InviteStore
}

func New(deps Deps) http.Handler {
	mux := http.NewServeMux()

	userH := handlers.NewUserHandler(deps.UserStore)
	groupH := handlers.NewGroupHandler(deps.GroupStore, deps.RoleStore)
	memberH := handlers.NewMemberHandler(deps.MemberStore, deps.RoleStore)
	rollH := handlers.NewRollHandler(deps.MemberStore, deps.RollStore, deps.RoleStore, deps.GroupStore)

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/me", userH.GetMe)

	protected.HandleFunc("POST /api/groups", groupH.Create)
	protected.HandleFunc("GET /api/groups", groupH.List)
	protected.HandleFunc("GET /api/groups/{id}", groupH.Get)
	protected.HandleFunc("PUT /api/groups/{id}", groupH.Update)
	protected.HandleFunc("DELETE /api/groups/{id}", groupH.Delete)

	protected.HandleFunc("GET /api/groups/{gid}/members", memberH.List)
	protected.HandleFunc("POST /api/groups/{gid}/members", memberH.Add)
	protected.HandleFunc("PUT /api/groups/{gid}/members/{id}", memberH.Update)
	protected.HandleFunc("DELETE /api/groups/{gid}/members/{id}", memberH.Delete)
	protected.HandleFunc("PATCH /api/groups/{gid}/members/{id}/points", memberH.AdjustPoints)

	protected.HandleFunc("POST /api/groups/{gid}/rolls", rollH.ExecuteRoll)
	protected.HandleFunc("GET /api/groups/{gid}/rolls", rollH.ListRolls)
	protected.HandleFunc("GET /api/groups/{gid}/rolls/stats", rollH.GetStats)

	authed := middleware.Auth(deps.Verifier)(protected)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/health" {
			handlers.HealthCheck(w, r)
			return
		}
		authed.ServeHTTP(w, r)
	})

	_ = mux // mux is not directly used; all routing is via handler chain

	return middleware.Logging(middleware.CORS(handler))
}
