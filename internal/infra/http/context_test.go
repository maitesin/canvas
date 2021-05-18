package http_test

import (
	"context"
	"testing"

	httpx "github.com/maitesin/sketch/internal/infra/http"
	//nolint: depguard
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLoggerContext(t *testing.T) {
	ctx := context.Background()
	logger := log.New()

	ctx = httpx.ContextWithLogger(ctx, logger)
	require.Equal(t, logger, httpx.LoggerFromContext(ctx))
}
