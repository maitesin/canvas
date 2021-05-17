package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/maitesin/sketch/config"
	"github.com/maitesin/sketch/internal/infra/ascii"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	sqlx "github.com/maitesin/sketch/internal/infra/sql"
	"github.com/upper/db/v4/adapter/postgresql"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Failed to generate configuration %s\n", err)
		return
	}

	dbConn, err := sql.Open("postgres", cfg.SQL.DatabaseURL())
	if err != nil {
		fmt.Printf("Failed to open connection to the DB: %s\n", err)
		return
	}
	defer dbConn.Close()

	pgConn, err := postgresql.New(dbConn)
	if err != nil {
		fmt.Printf("Failed to initialize connection with the DB: %s\n", err)
		return
	}
	defer pgConn.Close()

	canvasRepository := sqlx.NewCanvasRepository(pgConn)
	renderer := ascii.Renderer{}

	err = http.ListenAndServe(
		strings.Join([]string{cfg.HTTP.Host, cfg.HTTP.Port}, ":"),
		httpx.DefaultRouter(cfg.Canvas, canvasRepository, renderer),
	)
	if err != nil {
		fmt.Printf("Failed to start service: %s\n", err)
	}
}
