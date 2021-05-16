package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/maitesin/sketch/config"
	"github.com/maitesin/sketch/internal/infra/ascii"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	"github.com/maitesin/sketch/internal/infra/mem"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("failed to generate configuration %s\n", err)
		return
	}
	canvasRepository := &mem.CanvasRepository{}
	renderer := ascii.Renderer{}

	err = http.ListenAndServe(
		strings.Join([]string{cfg.HTTP.Host, cfg.HTTP.Port}, ":"),
		httpx.DefaultRouter(cfg.Canvas, canvasRepository, renderer))
	if err != nil {
		fmt.Printf("failed to start service: %s\n", err)
	}
}
