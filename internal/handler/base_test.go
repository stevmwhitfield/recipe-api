package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevmwhitfield/recipe-api/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.Root(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "root", w.Body.String())
}

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	handler.Ping(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPanic(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		handler.Panic(w, req)
	})
}
