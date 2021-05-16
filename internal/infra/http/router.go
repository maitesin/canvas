package http

import (
	httpx "net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/maitesin/sketch/internal/app"
)

func DefaultRouter(repository app.CanvasRepository, renderer Renderer) httpx.Handler {
	router := chi.NewRouter()
	router.Use(middleware.DefaultLogger)

	router.Post("/canvas", CreateCanvasHandler(app.NewCreateCanvasHandler(repository, 12, 32)))
	router.Get("/canvas/{canvasID}", RenderCanvasHandler(app.NewRetrieveCanvasHandler(repository), renderer))
	router.Post("/canvas/{canvasID}", AddTaskHandler(
		app.NewDrawRectangleHandler(repository),
		app.NewAddFillHandler(repository),
	))

	return router
}
