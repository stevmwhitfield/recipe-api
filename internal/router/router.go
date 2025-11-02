package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	"github.com/stevmwhitfield/recipe-api/internal/app"
	customMiddleware "github.com/stevmwhitfield/recipe-api/internal/middleware"
)

func InitRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://example.com"},
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Handlers
	r.Get("/", app.BaseHandler.Root)
	r.Get("/ping", app.BaseHandler.Ping)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(customMiddleware.APIVersionCtx("v1"))
		r.Mount("/recipes", app.RecipeHandler.Routes())
	})

	return r
}
