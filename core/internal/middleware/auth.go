package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/angelov-todor/lootforge/core/internal/auth"
	"github.com/angelov-todor/lootforge/core/internal/httputil"
)

type contextKey string

const userContextKey contextKey = "auth_user"

func Auth(verifier auth.TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				httputil.WriteError(w, 401, "missing or malformed authorization header")
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")

			user, err := verifier.VerifyToken(r.Context(), token)
			if err != nil {
				httputil.WriteError(w, 401, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) *auth.AuthUser {
	user, _ := ctx.Value(userContextKey).(*auth.AuthUser)
	return user
}

func UserIDFromContext(ctx context.Context) string {
	user := UserFromContext(ctx)
	if user == nil {
		return ""
	}
	return user.UID
}
