package http

import (
	"context"
	httpx "net/http"

	middleware "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi"
	"github.com/maitesin/sketch/internal/app"
	log "github.com/sirupsen/logrus" //nolint: depguard
)

func DefaultRouter(ctx context.Context, cfg app.Config, repository app.CanvasRepository, renderer Renderer) httpx.Handler {
	logger := LoggerFromContext(ctx)

	router := chi.NewRouter()
	router.Use(middleware.Logger("router", logger))

	router.Post("/canvas", loggerMiddleware(logger, CreateCanvasHandler(app.NewCreateCanvasHandler(repository, cfg.Height, cfg.Width))))
	router.Get("/canvas/{canvasID}", loggerMiddleware(logger, RenderCanvasHandler(app.NewRetrieveCanvasHandler(repository), renderer)))
	router.Post("/canvas/{canvasID}", loggerMiddleware(logger, AddTaskHandler(
		app.NewDrawRectangleHandler(repository),
		app.NewAddFillHandler(repository),
	)))

	return router
}

func loggerMiddleware(logger *log.Logger, next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(writer httpx.ResponseWriter, request *httpx.Request) {
		request = request.WithContext(ContextWithLogger(request.Context(), logger))
		next.ServeHTTP(writer, request)
	}
}
