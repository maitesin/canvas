package app_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
	"github.com/stretchr/testify/require"
)

func validCreateCanvasCmd() app.Command {
	return app.CreateCanvasCmd{
		ID: uuid.New(),
	}
}

type invalidCmd struct{}

func (i invalidCmd) Name() string { return "invalidCmd" }

func TestCreateCanvasHandler(t *testing.T) {
	tests := []struct {
		name              string
		command           app.Command
		repositoryMutator repositoryMutator
		expectedErr       error
	}{
		{
			name: `Given a valid command and a working canvas repository
                   when the create canvas handler is executed
                   then no error is returned`,
			command:           validCreateCanvasCmd(),
			repositoryMutator: noopRepositoryMutator,
		},
		{
			name: `Given an invalid command and a working canvas repository
                   when the create canvas handler is executed
                   then an invalid command error is returned`,
			command:           invalidCmd{},
			repositoryMutator: noopRepositoryMutator,
			expectedErr:       app.InvalidCommandError{},
		},
		{
			name: `Given a valid command and a non-working canvas repository
                   when the create canvas handler is executed
                   then an error returned`,
			command: validCreateCanvasCmd(),
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					InsertFunc: func(context.Context, domain.Canvas) error {
						return app.CanvasNotFound{}
					},
				}
				return repository
			},
			expectedErr: app.CanvasNotFound{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := tt.repositoryMutator(validCanvasRepository())
			handler := app.NewCreateCanvasHandler(repository)

			err := handler.Handle(context.Background(), tt.command)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func validDrawRectangleCmd() app.Command {
	return app.DrawRectangleCmd{
		CanvasID:    uuid.New(),
		RectangleID: uuid.New(),
		Point:       domain.NewPoint(uint(0), uint(0)),
		Height:      uint(10),
		Width:       uint(10),
		Filler:      '0',
		Outline:     'X',
	}
}

func TestDrawRectangleHandler(t *testing.T) {
	tests := []struct {
		name              string
		command           app.Command
		repositoryMutator repositoryMutator
		expectedErr       error
	}{
		{
			name: `Given a valid command and a working canvas repository
                   when the draw rectangle handler is executed
                   then no error is returned`,
			command:           validDrawRectangleCmd(),
			repositoryMutator: noopRepositoryMutator,
		},
		{
			name: `Given an invalid command and a working canvas repository
                   when the draw rectangle handler is executed
                   then an invalid command error returned`,
			command:           invalidCmd{},
			repositoryMutator: noopRepositoryMutator,
			expectedErr:       app.InvalidCommandError{},
		},
		{
			name: `Given a valid command and a canvas repository that fails the find by ID operation
                   when the draw rectangle handler is executed
                   then an error is returned`,
			command: validDrawRectangleCmd(),
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, app.CanvasNotFound{}
					},
				}
				return repository
			},
			expectedErr: app.CanvasNotFound{},
		},
		{
			name: `Given a valid command and a canvas repository that fails the update operation
                   when the draw rectangle handler is executed
                   then an error is returned`,
			command: validDrawRectangleCmd(),
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, nil
					},
					UpdateFunc: func(context.Context, domain.Canvas) error {
						return app.CanvasNotFound{}
					},
				}
				return repository
			},
			expectedErr: app.CanvasNotFound{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := tt.repositoryMutator(validCanvasRepository())
			handler := app.NewDrawRectangleHandler(repository)

			err := handler.Handle(context.Background(), tt.command)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func validAddFillCmd() app.Command {
	return app.AddFillCmd{
		CanvasID: uuid.New(),
		FillID:   uuid.New(),
		Point:    domain.NewPoint(uint(12), uint(12)),
		Filler:   '-',
	}
}

func TestAddFillHandler(t *testing.T) {
	tests := []struct {
		name              string
		command           app.Command
		repositoryMutator repositoryMutator
		expectedErr       error
	}{
		{
			name: `Given a valid command and a working canvas repository
                   when the add fill handler is executed
                   then no error is returned`,
			command:           validAddFillCmd(),
			repositoryMutator: noopRepositoryMutator,
		},
		{
			name: `Given an invalid command and a working canvas repository
                   when the add fill handler is executed
                   then an invalid command error returned`,
			command:           invalidCmd{},
			repositoryMutator: noopRepositoryMutator,
			expectedErr:       app.InvalidCommandError{},
		},
		{
			name: `Given a valid command and a canvas repository that fails the find by ID operation
                   when the add fill handler is executed
                   then an error is returned`,
			command: validAddFillCmd(),
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, app.CanvasNotFound{}
					},
				}
				return repository
			},
			expectedErr: app.CanvasNotFound{},
		},
		{
			name: `Given a valid command and a canvas repository that fails the update operation
                   when the add fill handler is executed
                   then an error is returned`,
			command: validAddFillCmd(),
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, nil
					},
					UpdateFunc: func(context.Context, domain.Canvas) error {
						return app.CanvasNotFound{}
					},
				}
				return repository
			},
			expectedErr: app.CanvasNotFound{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := tt.repositoryMutator(validCanvasRepository())
			handler := app.NewAddFillHandler(repository)

			err := handler.Handle(context.Background(), tt.command)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
