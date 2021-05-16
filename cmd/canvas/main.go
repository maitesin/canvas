package main

import (
	"fmt"
	"net/http"

	"github.com/maitesin/sketch/internal/infra/ascii"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	"github.com/maitesin/sketch/internal/infra/mem"
)

func main() {
	canvasRepository := &mem.CanvasRepository{}
	renderer := ascii.Renderer{}
	router := httpx.DefaultRouter(canvasRepository, renderer)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("failed to start service: %s\n", err)
	}
}
