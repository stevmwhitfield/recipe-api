package router_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevmwhitfield/recipe-api/internal/app"
	"github.com/stevmwhitfield/recipe-api/internal/handler"
	"github.com/stevmwhitfield/recipe-api/internal/router"
	"github.com/stretchr/testify/assert"
)

func TestRecipeRoutes(t *testing.T) {
	mockBaseHandler := &handler.BaseHandler{}
	mockRecipeHandler := &handler.RecipeHandler{}
	mockApp := &app.Application{
		Logger:        slog.New(slog.NewJSONHandler(io.Discard, nil)),
		BaseHandler:   mockBaseHandler,
		RecipeHandler: mockRecipeHandler,
	}

	r := router.InitRoutes(mockApp)

	tests := []struct {
		name       string
		method     string
		uri        string
		wantStatus int
	}{
		{"list recipes", http.MethodGet, "/api/v1/recipes", http.StatusOK},
		{"get recipe", http.MethodGet, "/api/v1/recipes/019a40de-02cd-7865-84ae-c038b75596f5", http.StatusOK},
		{"create recipe", http.MethodPost, "/api/v1/recipes", http.StatusCreated},
		{"update recipe", http.MethodPut, "/api/v1/recipes/019a40de-02cd-7865-84ae-c038b75596f5", http.StatusOK},
		{"delete recipe", http.MethodDelete, "/api/v1/recipes/019a40de-02cd-7865-84ae-c038b75596f5", http.StatusNoContent},
		{"route not found", http.MethodGet, "/api/v1/foobar", http.StatusNotFound},
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
