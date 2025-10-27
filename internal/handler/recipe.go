package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RecipeHandler struct{}

func NewRecipeHandler() *RecipeHandler {
	return &RecipeHandler{}
}

func (h *RecipeHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ListRecipes)
	r.Post("/", h.CreateRecipe)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetRecipe)
		r.Put("/", h.UpdateRecipe)
		r.Delete("/", h.DeleteRecipe)
	})

	return r
}

func (h *RecipeHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RecipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
