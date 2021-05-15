package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/domain"
	"github.com/stretchr/testify/require"
)

func validDrawRectangle() domain.DrawRectangle {
	return domain.NewDrawRectangle(
		uuid.New(),
		domain.NewPoint(0, 0),
		10,
		10,
		'0',
		'X',
		time.Now().UTC(),
	)
}

func validFill() domain.Fill {
	return domain.NewFill(
		uuid.New(),
		domain.NewPoint(12, 12),
		'-',
		time.Now().UTC(),
	)
}

func validCanvas() *domain.Canvas {
	rectangle := validDrawRectangle()
	fill := validFill()
	canvas := domain.NewCanvas(
		uuid.New(),
		30,
		30,
		[]domain.Task{rectangle, fill},
		time.Now().UTC(),
	)

	return &canvas
}

type canvasMutator func(*domain.Canvas) *domain.Canvas

var noopCanvasMutator = func(c *domain.Canvas) *domain.Canvas { return c }

func TestCanvas(t *testing.T) {
	tests := []struct {
		name                  string
		mutator               canvasMutator
		rectangles            []domain.DrawRectangle
		fills                 []domain.Fill
		expectedNumberOfTasks int
	}{
		{
			name: `Given a canvas containing a drawing rectangle and a fill operation,
                   when the Tasks method is called,
                   then returns the number of tasks of the canvas, 2 in this case`,
			mutator:               noopCanvasMutator,
			expectedNumberOfTasks: 2,
		},
		{

			name: `Given a canvas containing a drawing rectangle and a fill operation,
                   when the Tasks method is called after adding another drawing rectangle operation,
                   then returns the number of tasks of the canvas, 3 in this case`,
			mutator:               noopCanvasMutator,
			rectangles:            []domain.DrawRectangle{validDrawRectangle()},
			expectedNumberOfTasks: 3,
		},
		{

			name: `Given a canvas containing a drawing rectangle and a fill operation,
                   when the Tasks method is called after adding another fill operation,
                   then returns the number of tasks of the canvas, 3 in this case`,
			mutator:               noopCanvasMutator,
			fills:                 []domain.Fill{validFill()},
			expectedNumberOfTasks: 3,
		},
		{

			name: `Given a canvas containing a drawing rectangle and a fill operation,
                   when the Tasks method is called after adding another drawing rectangle and another fill operation,
                   then returns the number of tasks of the canvas, 4 in this case`,
			mutator:               noopCanvasMutator,
			rectangles:            []domain.DrawRectangle{validDrawRectangle()},
			fills:                 []domain.Fill{validFill()},
			expectedNumberOfTasks: 4,
		},
		{

			name: `Given a canvas containing no operations,
                   when the Tasks method is called,
                   then returns the number of tasks of the canvas, 0 in this case`,
			mutator: func(canvas *domain.Canvas) *domain.Canvas {
				c := domain.NewCanvas(
					uuid.New(),
					30,
					30,
					nil,
					time.Now().UTC(),
				)
				return &c
			},
			expectedNumberOfTasks: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := tt.mutator(validCanvas())

			for _, rectangle := range tt.rectangles {
				c.AddDrawRectangle(rectangle)
			}

			for _, fill := range tt.fills {
				c.AddFill(fill)
			}

			require.Len(t, c.Tasks(), tt.expectedNumberOfTasks)
		})
	}
}
