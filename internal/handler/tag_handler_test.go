package handler_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stevmwhitfield/recipe-api/internal/handler"
	"github.com/stevmwhitfield/recipe-api/internal/model"
	"github.com/stevmwhitfield/recipe-api/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//#region mocks

type MockTagStore struct {
	mock.Mock
}

func (m *MockTagStore) ListTags() ([]model.Tag, error) {
	args := m.Called()
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagStore) CreateTag(t *model.Tag) (*model.Tag, error) {
	args := m.Called(t)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

//#endregion

//#region tests

func TestTagHandler(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		uri       string
		data      io.Reader           // optional
		setupMock func(*MockTagStore) // optional
		wantCode  int
		wantBody  interface{}
	}{
		{
			name:   "list tags",
			method: http.MethodGet,
			uri:    "/",
			setupMock: func(m *MockTagStore) {
				m.On("ListTags").Return([]model.Tag{
					{ID: "1", Name: "t1"},
					{ID: "2", Name: "t2"},
				}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: util.Envelope{"tags": []model.Tag{{ID: "1", Name: "t1"}, {ID: "2", Name: "t2"}}, "total": 2},
		},
		{
			name:   "list tags with error",
			method: http.MethodGet,
			uri:    "/",
			setupMock: func(m *MockTagStore) {
				m.On("ListTags").Return([]model.Tag{}, errors.New("database error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: util.Envelope{"error": "failed to fetch tags"},
		},
		{
			name:   "create tag",
			method: http.MethodPost,
			uri:    "/",
			data:   strings.NewReader(`{ "name": "t1" }`),
			setupMock: func(m *MockTagStore) {
				m.On("CreateTag", mock.AnythingOfType("*model.Tag")).Return(&model.Tag{Name: "t1"}, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: &model.Tag{Name: "t1"},
		},
		{
			name:     "create tag with invalid json",
			method:   http.MethodPost,
			uri:      "/",
			data:     strings.NewReader(`{ "name": "t1" `),
			wantCode: http.StatusBadRequest,
			wantBody: util.Envelope{"error": "invalid request body"},
		},
		{
			name:     "create tag with missing name",
			method:   http.MethodPost,
			uri:      "/",
			data:     strings.NewReader(`{ "foo": "bar" }`),
			wantCode: http.StatusBadRequest,
			wantBody: util.Envelope{"error": "name cannot be blank"},
		},
		{
			name:     "create tag with blank name",
			method:   http.MethodPost,
			uri:      "/",
			data:     strings.NewReader(`{ "name": "" }`),
			wantCode: http.StatusBadRequest,
			wantBody: util.Envelope{"error": "name cannot be blank"},
		},
		{
			name:     "create tag with wrong name type",
			method:   http.MethodPost,
			uri:      "/",
			data:     strings.NewReader(`{ "name": 123 }`),
			wantCode: http.StatusBadRequest,
			wantBody: util.Envelope{"error": "invalid request body"},
		},
		{
			name:   "create tag with failure creating tag",
			method: http.MethodPost,
			uri:    "/",
			data:   strings.NewReader(`{ "name": "t1" }`),
			setupMock: func(m *MockTagStore) {
				m.On("CreateTag", mock.AnythingOfType("*model.Tag")).Return(nil, errors.New("database error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: util.Envelope{"error": "failed to create tag"},
		},
		{
			name:   "create tag with duplicate name",
			method: http.MethodPost,
			uri:    "/",
			data:   strings.NewReader(`{ "name": "t1" }`),
			setupMock: func(m *MockTagStore) {
				m.On("CreateTag", mock.AnythingOfType("*model.Tag")).Return(nil, errors.New("UNIQUE constraint failed: tags.name"))
			},
			wantCode: http.StatusConflict,
			wantBody: util.Envelope{"error": "tag with that name already exists"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

			mockStore := &MockTagStore{}
			if tt.setupMock != nil {
				tt.setupMock(mockStore)
			}

			h := handler.NewTagHandler(logger, mockStore)

			r := chi.NewRouter()
			r.Mount("/", h.Routes())

			req := httptest.NewRequest(tt.method, tt.uri, tt.data)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantBody != nil {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

				wantJSON, _ := json.Marshal(tt.wantBody)
				assert.JSONEq(t, string(wantJSON), w.Body.String())
			}

			mockStore.AssertExpectations(t)
		})
	}
}

//#endregion
