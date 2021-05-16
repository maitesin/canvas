package http

import (
	httpx "net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/maitesin/sketch/internal/app"
)

func DefaultRouter(cfg app.Config, repository app.CanvasRepository, renderer Renderer) httpx.Handler {
	router := chi.NewRouter()
	router.Use(middleware.DefaultLogger)

	router.Post("/canvas", CreateCanvasHandler(app.NewCreateCanvasHandler(repository, cfg.Height, cfg.Width)))
	router.Get("/canvas/{canvasID}", RenderCanvasHandler(app.NewRetrieveCanvasHandler(repository), renderer))
	router.Post("/canvas/{canvasID}", AddTaskHandler(
		app.NewDrawRectangleHandler(repository),
		app.NewAddFillHandler(repository),
	))

	return router
}
