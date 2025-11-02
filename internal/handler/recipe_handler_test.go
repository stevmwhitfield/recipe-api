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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stevmwhitfield/recipe-api/internal/handler"
	"github.com/stevmwhitfield/recipe-api/internal/model"
	"github.com/stevmwhitfield/recipe-api/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mocks

type MockRecipeStore struct {
	mock.Mock
}

func (m *MockRecipeStore) ListRecipes() ([]model.Recipe, error) {
	args := m.Called()
	return args.Get(0).([]model.Recipe), args.Error(1)
}
func (m *MockRecipeStore) CreateRecipe(r *model.Recipe) (*model.Recipe, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Recipe), args.Error(1)

}
func (m *MockRecipeStore) GetRecipeByID(id string) (*model.Recipe, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Recipe), args.Error(1)
}
func (m *MockRecipeStore) UpdateRecipe(r *model.Recipe) (*model.Recipe, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Recipe), args.Error(1)
}
func (m *MockRecipeStore) DeleteRecipe(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// tests

// TODO: finish writing tests for all handlers
func TestRecipeHandler(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		uri       string
		data      io.Reader // optional
		setupMock func(*MockRecipeStore)
		wantCode  int
		wantBody  util.Envelope // optional
	}{
		{
			name:   "list recipes",
			method: http.MethodGet,
			uri:    "/",
			setupMock: func(m *MockRecipeStore) {
				m.On("ListRecipes").Return(getListRecipeData(), nil)
			},
			wantCode: http.StatusOK,
			wantBody: util.Envelope{"recipes": getListRecipeData(), "total": 2},
		},
		{
			name:   "list recipes with error",
			method: http.MethodGet,
			uri:    "/",
			setupMock: func(m *MockRecipeStore) {
				m.On("ListRecipes").Return([]model.Recipe{}, errors.New("database error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: util.Envelope{"error": "failed to fetch recipes"},
		},
		{
			name:   "create recipe",
			method: http.MethodPost,
			uri:    "/",
			data:   getNewRecipeData(),
			setupMock: func(m *MockRecipeStore) {
				m.On("CreateRecipe", mock.AnythingOfType("*model.Recipe")).Return(
					&model.Recipe{
						Name:            "Classic Pancakes",
						Servings:        4,
						PrepTimeSeconds: 600,
						CookTimeSeconds: 900,
					}, nil,
				)
			},
			wantCode: http.StatusCreated,
			wantBody: util.Envelope{"recipe": &model.Recipe{
				Name:            "Classic Pancakes",
				Servings:        4,
				PrepTimeSeconds: 600,
				CookTimeSeconds: 900,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

			mockStore := &MockRecipeStore{}
			tt.setupMock(mockStore)

			h := handler.NewRecipeHandler(logger, mockStore)

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

func getListRecipeData() []model.Recipe {
	baseTime, _ := time.Parse(time.RFC3339, "2025-11-01T19:32:00Z")

	return []model.Recipe{
		{
			ID:              "019a40de-02cd-7865-84ae-c038b75596f5",
			Slug:            "classic-pancakes",
			Name:            "Classic Pancakes",
			Servings:        4,
			PrepTimeSeconds: 600,
			CookTimeSeconds: 900,
			Ingredients: []model.Ingredient{
				{ID: "i1", Name: "Flour", Quantity: 2, Unit: "cup", Note: "All-purpose"},
				{ID: "i2", Name: "Milk", Quantity: 1, Unit: "cup", Note: "Whole"},
				{ID: "i3", Name: "Eggs", Quantity: 2, Unit: "", Note: "Large"},
				{ID: "i4", Name: "Butter", Quantity: 2, Unit: "tbsp", Note: "Melted"},
			},
			Instructions: []model.Instruction{
				{ID: "s1", StepNumber: 1, Description: "Whisk together flour, milk, and eggs in a large bowl."},
				{ID: "s2", StepNumber: 2, Description: "Add melted butter and stir until smooth."},
				{ID: "s3", StepNumber: 3, Description: "Pour 1/4 cup batter onto a hot griddle and cook until golden."},
			},
			Tags: []model.Tag{
				{ID: "t1", Name: "Breakfast"},
				{ID: "t2", Name: "Easy"},
			},
			CreatedAt: baseTime.Add(-48 * time.Hour),
			UpdatedAt: baseTime.Add(-24 * time.Hour),
		},
		{
			ID:              "019a40de-02cd-7bc7-b171-710c99947f08",
			Slug:            "spaghetti-bolognese",
			Name:            "Spaghetti Bolognese",
			Servings:        6,
			PrepTimeSeconds: 900,
			CookTimeSeconds: 3600,
			Ingredients: []model.Ingredient{
				{ID: "i5", Name: "Spaghetti", Quantity: 500, Unit: "g", Note: ""},
				{ID: "i6", Name: "Ground Beef", Quantity: 500, Unit: "g", Note: "Lean"},
				{ID: "i7", Name: "Tomato Sauce", Quantity: 2, Unit: "cups", Note: ""},
				{ID: "i8", Name: "Onion", Quantity: 1, Unit: "", Note: "Chopped"},
				{ID: "i9", Name: "Garlic", Quantity: 2, Unit: "cloves", Note: "Minced"},
			},
			Instructions: []model.Instruction{
				{ID: "s4", StepNumber: 1, Description: "Cook spaghetti according to package directions."},
				{ID: "s5", StepNumber: 2, Description: "Brown beef with onion and garlic in a pan."},
				{ID: "s6", StepNumber: 3, Description: "Add tomato sauce and simmer for 30 minutes."},
				{ID: "s7", StepNumber: 4, Description: "Serve sauce over cooked spaghetti."},
			},
			Tags: []model.Tag{
				{ID: "t3", Name: "Dinner"},
				{ID: "t4", Name: "Italian"},
				{ID: "t5", Name: "Pasta"},
			},
			CreatedAt: baseTime.Add(-72 * time.Hour),
			UpdatedAt: baseTime,
		},
	}
}

func getNewRecipeData() io.Reader {
	js := `{
		"name": "Classic Pancakes",
		"servings": 4,
		"prepTimeSeconds": 600,
		"cookTimeSeconds": 900,
		"ingredients": [
			{"ingredientId": "i1", "quantity": 2, "unit": "cup", "note": "All-purpose"},
			{"ingredientId": "i2", "quantity": 1, "unit": "cup", "note": "Whole"},
			{"ingredientId": "i3", "quantity": 2, "unit": "", "note": "Large"},
			{"ingredientId": "i4", "quantity": 2, "unit": "tbsp", "note": "Melted"}
		],
		"instructions": [
			{"stepNumber": 1, "description": "Whisk together flour, milk, and eggs in a large bowl."},
			{"stepNumber": 2, "description": "Add melted butter and stir until smooth."},
			{"stepNumber": 3, "description": "Pour 1/4 cup batter onto a hot griddle and cook until golden."}
		],
		"tagIds": ["t1", "t2"]
	}`
	return strings.NewReader(js)
}
