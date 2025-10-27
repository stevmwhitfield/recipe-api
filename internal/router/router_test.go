package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevmwhitfield/recipe-api/internal/router"
	"github.com/stretchr/testify/assert"
)

func TestPanicRecovery(t *testing.T) {
	r := router.InitRoutes()

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRecipeRoutes(t *testing.T) {
	r := router.InitRoutes()

	tests := []struct {
		name       string
		method     string
		uri        string
		wantStatus int
	}{
		{"list recipes", http.MethodGet, "/api/v1/recipes", http.StatusNotImplemented},
		{"get recipe", http.MethodGet, "/api/v1/recipes/123", http.StatusNotImplemented},
		{"create recipe", http.MethodPost, "/api/v1/recipes", http.StatusNotImplemented},
		{"update recipe", http.MethodPut, "/api/v1/recipes/123", http.StatusNotImplemented},
		{"delete recipe", http.MethodDelete, "/api/v1/recipes/123", http.StatusNotImplemented},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.uri, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
