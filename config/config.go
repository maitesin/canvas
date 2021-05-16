package config

import (
	"os"
	"strconv"

	"github.com/maitesin/sketch/internal/app"
	httpx "github.com/maitesin/sketch/internal/infra/http"
)

const (
	defaultCanvasHeight = "12"
	defaultCanvasWidth  = "32"
	defaultPort         = "8080"
	defaultHost         = ""
)

// Config defines the general configuration of the service
type Config struct {
	HTTP   httpx.Config
	Canvas app.Config
}

func New() (Config, error) {
	canvasHeight, err := strconv.Atoi(getEnvOrDefault("CANVAS_HEIGHT", defaultCanvasHeight))
	if err != nil {
		return Config{}, err
	}

	canvasWidth, err := strconv.Atoi(getEnvOrDefault("CANVAS_WIDTH", defaultCanvasWidth))
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTP: httpx.Config{
			Host: getEnvOrDefault("HOST", defaultHost),
			Port: getEnvOrDefault("PORT", defaultPort),
		},
		Canvas: app.Config{
			Height: canvasHeight,
			Width:  canvasWidth,
		},
	}, nil
}

func getEnvOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}

	return defaultValue
}
