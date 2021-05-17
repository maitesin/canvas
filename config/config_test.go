package config_test

import (
	"os"
	"testing"

	"github.com/maitesin/sketch/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// unset environment variables
	variables := []string{
		"CANVAS_HEIGHT",
		"CANVAS_WIDTH",
		"HOST",
		"PORT",
		"DB_URL",
		"DB_SSL_MODE",
		"DB_BINARY_PARAMETERS",
	}
	for _, variable := range variables {
		err := os.Unsetenv(variable)
		require.NoError(t, err)
	}

	cfg, err := config.New()
	require.NoError(t, err)

	require.Equal(t, 12, cfg.Canvas.Height)
	require.Equal(t, 32, cfg.Canvas.Width)
	require.Equal(t, "", cfg.HTTP.Host)
	require.Equal(t, "8080", cfg.HTTP.Port)
	require.Equal(t, "postgres://postgres:postgres@localhost:54321/sketch", cfg.SQL.URL)
	require.Equal(t, "disable", cfg.SQL.SSLMode)
	require.Equal(t, "yes", cfg.SQL.BinaryParams)

	// set canvas height with not a number
	err = os.Setenv("CANVAS_HEIGHT", "nine")
	require.NoError(t, err)

	cfg, err = config.New()
	require.NotNil(t, err)

	err = os.Unsetenv("CANVAS_HEIGHT")
	require.NoError(t, err)

	// set canvas width with not a number
	err = os.Setenv("CANVAS_WIDTH", "twelve")
	require.NoError(t, err)

	cfg, err = config.New()
	require.NotNil(t, err)

	err = os.Unsetenv("CANVAS_WIDTH")
	require.NoError(t, err)

	// check that all the environment variables are being used correctly
	namesAndValues := [][2]string{
		{
			"CANVAS_HEIGHT", "100",
		},
		{
			"CANVAS_WIDTH", "200",
		},
		{
			"HOST", "Another",
		},
		{
			"PORT", "one",
		},
		{
			"DB_URL", "bites",
		},
		{
			"DB_SSL_MODE", "the",
		},
		{
			"DB_BINARY_PARAMETERS", "dust",
		},
	}

	for _, nameAndValue := range namesAndValues {
		err = os.Setenv(nameAndValue[0], nameAndValue[1])
		require.NoError(t, err)
	}

	cfg, err = config.New()
	require.NoError(t, err)

	require.Equal(t, 100, cfg.Canvas.Height)
	require.Equal(t, 200, cfg.Canvas.Width)
	require.Equal(t, "Another", cfg.HTTP.Host)
	require.Equal(t, "one", cfg.HTTP.Port)
	require.Equal(t, "bites", cfg.SQL.URL)
	require.Equal(t, "the", cfg.SQL.SSLMode)
	require.Equal(t, "dust", cfg.SQL.BinaryParams)
}
