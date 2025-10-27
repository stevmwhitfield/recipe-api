package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevmwhitfield/recipe-api/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAPIVersionCtx(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(middleware.APIVersionKey)
		if v != nil {
			w.Write([]byte(v.(string)))
		}
	})

	handler := middleware.APIVersionCtx("v1")(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, "v1", w.Body.String())
}
