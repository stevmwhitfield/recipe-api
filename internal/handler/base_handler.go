package handler

import (
	"log/slog"
	"net/http"

	utils "github.com/stevmwhitfield/recipe-api/internal/util"
)

type BaseHandler struct {
	logger *slog.Logger
}

func NewBaseHandler(l *slog.Logger) *BaseHandler {
	return &BaseHandler{
		logger: l,
	}
}

func (h *BaseHandler) Root(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "ok", "version": "v1"})
}

func (h *BaseHandler) Ping(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "pong"})
}
