// +build integration

package sql_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
	sqlx "github.com/maitesin/sketch/internal/infra/sql"
	"github.com/stretchr/testify/require"
)

func validCanvas() domain.Canvas {
	return domain.NewCanvas(
		uuid.New(),
		30,
		30,
		[]domain.Task{
			domain.NewDrawRectangle(
				uuid.New(),
				domain.NewPoint(0, 0),
				10,
				10,
				'0',
				'X',
				time.Now().UTC(),
			),
			domain.NewFill(
				uuid.New(),
				domain.NewPoint(12, 12),
				'-',
				time.Now().UTC(),
			),
		},
		time.Now().UTC(),
	)
}

func TestCanvasRepository_Insert(t *testing.T) {
	sess, dbCloser := setup(t)
	defer dbCloser()

	canvas := validCanvas()

	tests := []struct {
		name        string
		canvas      domain.Canvas
		fixtureID   uuid.UUID
		expectedErr error
	}{
		{
			name: `Given a working canvas repository and a valid canvas,
                   when the insert method is called,
                   then no error is returned.`,
			canvas:    validCanvas(),
			fixtureID: uuid.New(),
		},
		{
			name: `Given a working canvas repository and an invalid canvas,
                   when the insert method is called,
                   then an error is returned.`,
			canvas:      domain.NewCanvas(uuid.New(), 10, 10, []domain.Task{"I am invalid"}, time.Now().UTC()),
			fixtureID:   uuid.New(),
			expectedErr: errors.New(""),
		},
		{
			name: `Given a working canvas repository and a valid canvas that is already present in the DB,
                   when the insert method is called,
                   then no error is returned.`,
			canvas:    canvas,
			fixtureID: canvas.ID(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtures(t, sess, tt.fixtureID, uuid.New(), uuid.New())

			repository := sqlx.NewCanvasRepository(sess)

			err := repository.Insert(context.Background(), tt.canvas)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCanvasRepository_Update(t *testing.T) {
	sess, dbCloser := setup(t)
	defer dbCloser()

	tests := []struct {
		name        string
		canvas      domain.Canvas
		expectedErr error
	}{
		{
			name: `Given a working canvas repository and a valid canvas,
                   when the update method is called,
                   then no error is returned.`,
			canvas: validCanvas(),
		},
		{
			name: `Given a working canvas repository and an invalid canvas,
                   when the update method is called,
                   then an error is returned.`,
			canvas:      domain.NewCanvas(uuid.New(), 10, 10, []domain.Task{"I am invalid"}, time.Now().UTC()),
			expectedErr: errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtures(t, sess, tt.canvas.ID(), uuid.New(), uuid.New())

			repository := sqlx.NewCanvasRepository(sess)

			err := repository.Update(context.Background(), tt.canvas)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCanvasRepository_FindByID(t *testing.T) {
	sess, dbCloser := setup(t)
	defer dbCloser()

	canvas := validCanvas()

	tests := []struct {
		name        string
		canvas      domain.Canvas
		fixtureID   uuid.UUID
		expectedErr error
	}{
		{
			name: `Given a working canvas repository and a valid canvas that already exists in the DB,
                   when the find by id method is called,
                   then no error is returned.`,
			canvas:    canvas,
			fixtureID: canvas.ID(),
		},
		{
			name: `Given a working canvas repository and a valid canvas that does not exist in the DB,
                   when the find by id method is called,
                   then a canvas not found error is returned.`,
			canvas:      validCanvas(),
			fixtureID:   uuid.New(),
			expectedErr: app.CanvasNotFound{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtures(t, sess, tt.fixtureID, uuid.New(), uuid.New())

			repository := sqlx.NewCanvasRepository(sess)

			_, err := repository.FindByID(context.Background(), tt.canvas.ID())
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
