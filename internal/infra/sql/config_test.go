package sql_test

import (
	"testing"

	sqlx "github.com/maitesin/sketch/internal/infra/sql"
	"github.com/stretchr/testify/require"
)

func TestConfig_DatabaseURL(t *testing.T) {
	t.Parallel()

	cfg := sqlx.Config{
		URL:          "postgres://admin:admin@localhost:5432/db",
		SSLMode:      "disabled",
		BinaryParams: "yes",
	}

	require.Contains(t, cfg.DatabaseURL(), cfg.URL)
	require.Contains(t, cfg.DatabaseURL(), cfg.SSLMode)
	require.Contains(t, cfg.DatabaseURL(), cfg.BinaryParams)
}
