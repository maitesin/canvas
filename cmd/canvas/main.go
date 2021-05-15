package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/infra/ascii"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	"github.com/maitesin/sketch/internal/infra/mem"
)

func main() {
	canvasRepository := &mem.CanvasRepository{}
	renderer := ascii.Renderer{}

	router := chi.NewRouter()
	router.Use(middleware.DefaultLogger)

	router.Post("/canvas", httpx.CreateCanvasHandler(app.NewCreateCanvasHandler(canvasRepository, 12, 32)))
	router.Get("/canvas/{canvasID}", httpx.RenderCanvasHandler(app.NewRetrieveCanvasHandler(canvasRepository), renderer))
	router.Put("/canvas/{canvasID}", httpx.AddTaskHandler(
		app.NewDrawRectangleHandler(canvasRepository),
		app.NewAddFillHandler(canvasRepository),
	))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("failed to start service: %s\n", err)
	}
}
