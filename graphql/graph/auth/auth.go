package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey struct {
	name string
}

var tokenCtxKey = &contextKey{"token"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			words := strings.Fields(header)

			if len(words) != 2 {
				next.ServeHTTP(w, r)
				return
			}

			token := words[1]
			ctx := context.WithValue(r.Context(), tokenCtxKey, token)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) string {
	raw, _ := ctx.Value(tokenCtxKey).(string)
	return raw
}
