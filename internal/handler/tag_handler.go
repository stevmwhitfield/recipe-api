package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/stevmwhitfield/recipe-api/internal/model"
	"github.com/stevmwhitfield/recipe-api/internal/store"
	"github.com/stevmwhitfield/recipe-api/internal/util"
)

type TagHandler struct {
	logger   *slog.Logger
	tagStore store.TagStore
}

func NewTagHandler(l *slog.Logger, ts store.TagStore) *TagHandler {
	return &TagHandler{logger: l, tagStore: ts}
}

func (th *TagHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", th.ListTags)
	r.Post("/", th.CreateTag)

	return r
}

func (th *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := th.tagStore.ListTags()
	if err != nil {
		th.logger.Error("ListTags", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to fetch tags"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"tags": tags, "total": len(tags)})
}

func (th *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tag model.Tag

	err := json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		th.logger.Error("CreateTag", "error", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request body"})
		return
	}

	if tag.Name == "" {
		th.logger.Error("CreateTag", "error", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "name cannot be blank"})
		return
	}

	id, err := util.GenerateUUID()
	if err != nil {
		th.logger.Error("CreateTag", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to generate uuid"})
		return
	}
	tag.ID = id

	createdTag, err := th.tagStore.CreateTag(&tag)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			th.logger.Error("CreateTag", "error", err)
			util.WriteJSON(w, http.StatusConflict, util.Envelope{"error": "tag with that name already exists"})
			return
		}
		th.logger.Error("CreateTag", "error", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to create tag"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdTag)
}
