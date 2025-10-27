package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const APIVersionKey contextKey = "api.version"

func APIVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), APIVersionKey, version))
			next.ServeHTTP(w, r)
		})
	}
}
