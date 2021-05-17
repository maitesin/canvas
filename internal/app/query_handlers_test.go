package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
	"github.com/stretchr/testify/require"
)

func validCanvasRepository() app.CanvasRepository {
	return &CanvasRepositoryMock{
		FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
			return domain.NewCanvas(uuid.New(), 30, 30, nil, time.Now().UTC()), nil
		},
		InsertFunc: func(context.Context, domain.Canvas) error {
			return nil
		},
		UpdateFunc: func(context.Context, domain.Canvas) error {
			return nil
		},
	}
}

type repositoryMutator func(app.CanvasRepository) app.CanvasRepository

var noopRepositoryMutator = func(repository app.CanvasRepository) app.CanvasRepository { return repository }

func validQuery() app.Query {
	return app.RetrieveCanvasQuery{
		ID: uuid.New(),
	}
}

type invalidQuery struct{}

func (i invalidQuery) Name() string { return "invalidQuery" }

func TestRetrieveCanvasHandler(t *testing.T) {
	tests := []struct {
		name              string
		query             app.Query
		repositoryMutator repositoryMutator
		expectedErr       error
	}{
		{
			name: `Given a valid query and a working canvas repository
                   when the retrieve canvas query is handled
                   then a canvas is returned and no error is returned`,
			query:             validQuery(),
			repositoryMutator: noopRepositoryMutator,
		},
		{
			name: `Given an invalid query and a working canvas repository
                   when the retrieve canvas query is handled
                   then an invalid query error is returned`,
			query:             invalidQuery{},
			repositoryMutator: noopRepositoryMutator,
			expectedErr:       app.InvalidQueryError{},
		},
		{
			name: `Given a valid query and a non-working canvas repository
                   when the retrieve canvas query is handled
                   then an error is returned`,
			query: validQuery(),
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, app.CanvasNotFound{}
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
			handler := app.NewRetrieveCanvasHandler(repository)

			response, err := handler.Handle(context.Background(), tt.query)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
				_, ok := response.(domain.Canvas)
				require.True(t, ok)
			}
		})
	}
}
