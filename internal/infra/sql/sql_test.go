// +build integration

package sql_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/config"
	"github.com/maitesin/sketch/internal/domain"
	sqlx "github.com/maitesin/sketch/internal/infra/sql"
	"github.com/stretchr/testify/require"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

const migrationsDir = "../../../devops/db/migrations"

func TestMain(m *testing.M) {
	os.Setenv("DB_URL", "postgres://postgres:postgres@localhost:54321/sketch_test")
	os.Exit(m.Run())
}

type closer func()

func createDB(t *testing.T) {
	t.Helper()

	originalURL := os.Getenv("DB_URL")
	lastSlash := strings.LastIndex(originalURL, "/")
	dbURLWithoutDBName := originalURL[:lastSlash]
	dbName := originalURL[lastSlash+1:]

	err := os.Setenv("DB_URL", dbURLWithoutDBName)
	require.NoError(t, err)

	cfg, err := config.New()
	require.NoError(t, err)

	dbConn, err := sql.Open("postgres", cfg.SQL.DatabaseURL())
	require.NoError(t, err)

	queries := []string{
		fmt.Sprintf(`SELECT
		pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s' AND
			pid <> pg_backend_pid()`, dbName),
		fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, dbName),
		fmt.Sprintf("CREATE DATABASE %s", dbName),
	}
	for _, q := range queries {
		_, err := dbConn.Exec(q)
		require.NoError(t, err)
	}

	err = os.Setenv("DB_URL", originalURL)
	require.NoError(t, err)
}

func setup(t *testing.T) (db.Session, closer) {
	t.Helper()

	createDB(t)

	cfg, err := config.New()
	require.NoError(t, err)

	dbConn, err := sql.Open("postgres", cfg.SQL.DatabaseURL())
	require.NoError(t, err)

	pgConn, err := postgresql.New(dbConn)
	require.NoError(t, err)

	dbCloser := func() {
		defer dbConn.Close()
		defer pgConn.Close()
	}

	migrations, err := ioutil.ReadDir(migrationsDir)
	require.NoError(t, err)

	for _, migration := range migrations {
		content, err := ioutil.ReadFile(path.Join(migrationsDir, migration.Name()))
		require.NoError(t, err)
		_, err = pgConn.SQL().Exec(string(content))
		require.NoError(t, err)
	}

	return pgConn, dbCloser
}

func fixtures(t *testing.T, session db.Session, canvasID, rectangleID, fillID uuid.UUID) {
	t.Helper()

	rectangle := domain.NewDrawRectangle(
		rectangleID,
		domain.NewPoint(0, 0),
		10,
		10,
		'0',
		'X',
		time.Now().UTC(),
	)

	fill := domain.NewFill(
		fillID,
		domain.NewPoint(2, 2),
		'-',
		time.Now().UTC(),
	)

	canvas := domain.NewCanvas(
		canvasID,
		10,
		10,
		[]domain.Task{rectangle, fill},
		time.Now().UTC(),
	)

	repository := sqlx.NewCanvasRepository(session)
	err := repository.Insert(context.Background(), canvas) // To insert the canvas
	require.NoError(t, err)
	err = repository.Update(context.Background(), canvas) // To insert the tasks
	require.NoError(t, err)
}
