package main

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/maitesin/sketch/config"
	"github.com/maitesin/sketch/internal/infra/ascii"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	sqlx "github.com/maitesin/sketch/internal/infra/sql"
	log "github.com/sirupsen/logrus" //nolint: depguard
	"github.com/upper/db/v4/adapter/postgresql"
)

func main() {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	ctx := httpx.ContextWithLogger(context.Background(), logger)

	cfg, err := config.New()
	if err != nil {
		logger.Infof("Failed to generate configuration %s", err)
		return
	}

	dbConn, err := sql.Open("postgres", cfg.SQL.DatabaseURL())
	if err != nil {
		logger.Infof("Failed to open connection to the DB: %s\n", err)
		return
	}
	defer dbConn.Close()

	pgConn, err := postgresql.New(dbConn)
	if err != nil {
		logger.Infof("Failed to initialize connection with the DB: %s\n", err)
		return
	}
	defer pgConn.Close()

	canvasRepository := sqlx.NewCanvasRepository(pgConn)
	renderer := ascii.Renderer{}

	err = http.ListenAndServe(
		strings.Join([]string{cfg.HTTP.Host, cfg.HTTP.Port}, ":"),
		httpx.DefaultRouter(ctx, cfg.Canvas, canvasRepository, renderer),
	)
	if err != nil {
		logger.Infof("Failed to start service: %s\n", err)
	}
}
