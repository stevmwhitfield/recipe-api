package app

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/stevmwhitfield/recipe-api/internal/data/migrations"
	"github.com/stevmwhitfield/recipe-api/internal/handler"
	"github.com/stevmwhitfield/recipe-api/internal/store"
)

type Application struct {
	Logger        *slog.Logger
	BaseHandler   *handler.BaseHandler
	RecipeHandler *handler.RecipeHandler
	DB            *sql.DB
}

func NewApplication() (*Application, error) {
	// Database
	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Stores
	recipeStore := store.NewSQLiteRecipeStore(db)

	// Handlers
	baseHandler := handler.NewBaseHandler(logger)
	recipeHandler := handler.NewRecipeHandler(logger, recipeStore)

	app := &Application{
		Logger:        logger,
		BaseHandler:   baseHandler,
		RecipeHandler: recipeHandler,
		DB:            db,
	}

	return app, nil
}
