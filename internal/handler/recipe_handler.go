package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/stevmwhitfield/recipe-api/internal/store"
	"github.com/stevmwhitfield/recipe-api/internal/util"
)

type RecipeHandler struct {
	logger      *slog.Logger
	recipeStore store.RecipeStore
}

func NewRecipeHandler(l *slog.Logger, rs store.RecipeStore) *RecipeHandler {
	return &RecipeHandler{
		logger:      l,
		recipeStore: rs,
	}
}

func (h *RecipeHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ListRecipes)
	r.Post("/", h.CreateRecipe)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetRecipeByID)
		r.Put("/", h.UpdateRecipe)
		r.Delete("/", h.DeleteRecipe)
	})

	return r
}

func (h *RecipeHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.recipeStore.ListRecipes()
	if err != nil {
		h.logger.Error("ListRecipes", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to fetch recipes"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"recipes": recipes, "total": len(recipes)})
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe store.Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		h.logger.Error("CreateRecipe", "error", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request body"})
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		h.logger.Error("CreateRecipe", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to generate uuid"})
		return
	}
	recipe.ID = id.String()

	recipe.Slug = slug.Make(recipe.Name)

	createdRecipe, err := h.recipeStore.CreateRecipe(&recipe)
	if err != nil {
		h.logger.Error("CreateRecipe", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to create recipe"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"recipe": createdRecipe})
}

func (h *RecipeHandler) GetRecipeByID(w http.ResponseWriter, r *http.Request) {
	recipeID, err := util.ReadIDParam(r)
	if err != nil {
		h.logger.Error("GetRecipeByID", "error", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid recipe id"})
		return
	}

	recipe, err := h.recipeStore.GetRecipeByID(recipeID)
	if err != nil {
		h.logger.Error("GetRecipeByID", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to fetch recipe"})
		return
	}
	if recipe == nil {
		http.NotFound(w, r)
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"recipe": recipe})
}

func (h *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID, err := util.ReadIDParam(r)
	if err != nil {
		h.logger.Error("UpdateRecipe", "error", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid recipe id"})
		return
	}

	existingRecipe, err := h.recipeStore.GetRecipeByID(recipeID)
	if err != nil {
		h.logger.Error("UpdateRecipe", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to fetch recipe"})
		return
	}
	if existingRecipe == nil {
		http.NotFound(w, r)
		return
	}

	var recipeUpdateRequest struct {
		Name            *string             `json:"name"`
		Servings        *int                `json:"servings"`
		PrepTimeSeconds *int                `json:"prepTimeSeconds"`
		CookTimeSeconds *int                `json:"cookTimeSeconds"`
		Ingredients     []store.Ingredient  `json:"ingredients"`
		Instructions    []store.Instruction `json:"instructions"`
		Tags            []store.Tag         `json:"tags"`
	}

	err = json.NewDecoder(r.Body).Decode(&recipeUpdateRequest)
	if err != nil {
		h.logger.Error("UpdateRecipe", "error", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request body"})
		return
	}

	if recipeUpdateRequest.Name != nil {
		existingRecipe.Name = *recipeUpdateRequest.Name
	}
	if recipeUpdateRequest.Servings != nil {
		existingRecipe.Servings = *recipeUpdateRequest.Servings
	}
	if recipeUpdateRequest.PrepTimeSeconds != nil {
		existingRecipe.PrepTimeSeconds = *recipeUpdateRequest.PrepTimeSeconds
	}
	if recipeUpdateRequest.CookTimeSeconds != nil {
		existingRecipe.CookTimeSeconds = *recipeUpdateRequest.CookTimeSeconds
	}
	if recipeUpdateRequest.Ingredients != nil {
		existingRecipe.Ingredients = recipeUpdateRequest.Ingredients
	}
	if recipeUpdateRequest.Instructions != nil {
		existingRecipe.Instructions = recipeUpdateRequest.Instructions
	}
	if recipeUpdateRequest.Tags != nil {
		existingRecipe.Tags = recipeUpdateRequest.Tags
	}

	updatedRecipe, err := h.recipeStore.UpdateRecipe(existingRecipe)
	if err != nil {
		h.logger.Error("UpdateRecipe", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to update recipe"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"recipe": updatedRecipe})
}

func (h *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID, err := util.ReadIDParam(r)
	if err != nil {
		h.logger.Error("DeleteRecipe", "error", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid recipe id"})
		return
	}

	err = h.recipeStore.DeleteRecipe(recipeID)
	if err == sql.ErrNoRows {
		http.Error(w, "recipe not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.logger.Error("DeleteRecipe", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to delete recipe"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
