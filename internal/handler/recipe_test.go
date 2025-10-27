package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevmwhitfield/recipe-api/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestRecipeHandler(t *testing.T) {
	h := handler.NewRecipeHandler()

	t.Run("list recipes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		h.ListRecipes(w, req)

		assert.Equal(t, http.StatusNotImplemented, w.Code)
	})

	t.Run("get recipe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/123", nil)
		w := httptest.NewRecorder()

		h.GetRecipe(w, req)

		assert.Equal(t, http.StatusNotImplemented, w.Code)
	})

	t.Run("create recipe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()

		h.CreateRecipe(w, req)

		assert.Equal(t, http.StatusNotImplemented, w.Code)
	})

	t.Run("update recipe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/123", nil)
		w := httptest.NewRecorder()

		h.UpdateRecipe(w, req)

		assert.Equal(t, http.StatusNotImplemented, w.Code)
	})

	t.Run("delete recipe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/123", nil)
		w := httptest.NewRecorder()

		h.DeleteRecipe(w, req)

		assert.Equal(t, http.StatusNotImplemented, w.Code)
	})
}
