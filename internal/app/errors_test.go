package app_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/stretchr/testify/require"
)

func TestCanvasNotFound_Error(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	err := app.CanvasNotFound{ID: id}

	require.Contains(t, err.Error(), id.String())
}

func TestInvalidCommandError_Error(t *testing.T) {
	t.Parallel()

	cmd1 := app.CreateCanvasCmd{}
	cmd2 := app.DrawRectangleCmd{}
	err := app.InvalidCommandError{Received: cmd1, Expected: cmd2}

	require.Contains(t, err.Error(), cmd1.Name())
	require.Contains(t, err.Error(), cmd2.Name())

	cmd3 := app.AddFillCmd{}
	require.NotContains(t, err.Error(), cmd3.Name())
}

func TestInvalidQueryError_Error(t *testing.T) {
	t.Parallel()

	query1 := app.RetrieveCanvasQuery{}
	query2 := &QueryMock{
		NameFunc: func() string {
			return "wololo"
		},
	}

	err := app.InvalidQueryError{Expected: query1, Received: query2}
	require.Contains(t, err.Error(), query1.Name())
	require.Contains(t, err.Error(), query2.Name())
}
