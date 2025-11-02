package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevmwhitfield/recipe-api/internal/handler"
	"github.com/stevmwhitfield/recipe-api/internal/middleware"
	"github.com/stevmwhitfield/recipe-api/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestBaseHandler(t *testing.T) {
	h := handler.NewBaseHandler(nil)

	t.Run("root", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), middleware.APIVersionKey, "v1")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.Root(w, req)

		var res util.Envelope
		json.Unmarshal(w.Body.Bytes(), &res)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, util.Envelope{"status": "ok", "version": "v1"}, res)
	})

	t.Run("ping", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()

		h.Ping(w, req)

		var res util.Envelope
		json.Unmarshal(w.Body.Bytes(), &res)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, util.Envelope{"message": "pong"}, res)
	})
}
