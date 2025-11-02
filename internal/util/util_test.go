package util_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stevmwhitfield/recipe-api/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := util.Envelope{"message": "test"}

	err := util.WriteJSON(w, http.StatusOK, data)

	var body util.Envelope
	json.Unmarshal(w.Body.Bytes(), &body)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, data, body)
}

func TestReadIDParam_ValidID(t *testing.T) {
	rawID, err := uuid.NewV7()
	if err != nil {
		t.Errorf("failed to generate uuid, got: '%s', error: '%v'", rawID, err)
	}
	expectedID := rawID.String()

	r := httptest.NewRequest(http.MethodGet, "/recipes/"+expectedID, nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", expectedID)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	id, err := util.ReadIDParam(r)
	if id != expectedID {
		t.Errorf("expected: '%s', got: '%s', error: '%v'", expectedID, id, err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestReadIDParam_InvalidID(t *testing.T) {
	expectedID := "123"

	r := httptest.NewRequest(http.MethodGet, "/recipes/"+expectedID, nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", expectedID)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	id, err := util.ReadIDParam(r)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid id parameter type")
	assert.Equal(t, "", id)
}

func TestReadIDParam_MissingID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/recipes/", nil)
	ctx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	id, err := util.ReadIDParam(r)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid id parameter")
	assert.Equal(t, "", id)
}
