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
	canvas := domain.NewCanvas(
		uuid.New(),
		30,
		30,
		[]domain.Task{validDrawRectangle(), validFill()},
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
				err := c.AddDrawRectangle(rectangle)
				require.NoError(t, err)
			}

			for _, fill := range tt.fills {
				err := c.AddFill(fill)
				require.NoError(t, err)
			}

			require.Len(t, c.Tasks(), tt.expectedNumberOfTasks)
		})
	}
}

func TestCanvasGetters(t *testing.T) {
	t.Parallel()

	canvasID := uuid.New()
	canvasTime := time.Now().UTC()
	tasks := []domain.Task{validFill()}
	canvas := domain.NewCanvas(
		canvasID,
		30,
		10,
		tasks,
		canvasTime,
	)

	require.Equal(t, canvasID, canvas.ID())
	require.Equal(t, 30, canvas.Height())
	require.Equal(t, 10, canvas.Width())
	require.Equal(t, tasks, canvas.Tasks())
	require.Equal(t, canvasTime, canvas.CreatedAt())
}

func TestDrawRectangleGetters(t *testing.T) {
	t.Parallel()

	rectangleID := uuid.New()
	rectangleTime := time.Now().UTC()
	rectangle := domain.NewDrawRectangle(
		rectangleID,
		domain.NewPoint(4, 5),
		30,
		10,
		'X',
		'O',
		rectangleTime,
	)

	require.Equal(t, rectangleID, rectangle.ID())
	require.Equal(t, 4, rectangle.Point().X())
	require.Equal(t, 5, rectangle.Point().Y())
	require.Equal(t, 30, rectangle.Height())
	require.Equal(t, 10, rectangle.Width())
	require.Equal(t, 'X', rectangle.Filler())
	require.Equal(t, 'O', rectangle.Outline())
	require.Equal(t, rectangleTime, rectangle.CreatedAt())
}

func TestFillGetters(t *testing.T) {
	t.Parallel()

	fillID := uuid.New()
	fillTime := time.Now().UTC()
	fill := domain.NewFill(
		fillID,
		domain.NewPoint(4, 5),
		'-',
		fillTime,
	)

	require.Equal(t, fillID, fill.ID())
	require.Equal(t, 4, fill.Point().X())
	require.Equal(t, 5, fill.Point().Y())
	require.Equal(t, '-', fill.Filler())
	require.Equal(t, fillTime, fill.CreatedAt())
}

func TestCanvas_AddDrawRectangle(t *testing.T) {
	tests := []struct {
		name        string
		rectangle   domain.DrawRectangle
		expectedErr error
	}{
		{
			name: `Given a canvas with dimensions 30x30 and a draw rectangle with dimensions 30x30 from point 0x0,
                   when the AddDrawRectangle method is called,
                   then no error is returned`,
			rectangle: domain.NewDrawRectangle(uuid.New(), domain.NewPoint(0, 0), 30, 30, ' ', ' ', time.Now().UTC()),
		},
		{
			name: `Given a canvas with dimensions 30x30 and a draw rectangle with dimensions 30x31 from point 0x0,
                   when the AddDrawRectangle method is called,
                   then an out of bounds error is returned`,
			rectangle:   domain.NewDrawRectangle(uuid.New(), domain.NewPoint(0, 0), 31, 30, ' ', ' ', time.Now().UTC()),
			expectedErr: domain.ErrOutOfBounds,
		},
		{
			name: `Given a canvas with dimensions 30x30 and a draw rectangle with dimensions 31x30 from point 0x0,
                   when the AddDrawRectangle method is called,
                   then an out of bounds error is returned`,
			rectangle:   domain.NewDrawRectangle(uuid.New(), domain.NewPoint(0, 0), 30, 31, ' ', ' ', time.Now().UTC()),
			expectedErr: domain.ErrOutOfBounds,
		},
		{
			name: `Given a canvas with dimensions 30x30 and a draw rectangle with dimensions 30x30 from point 1x0,
                   when the AddDrawRectangle method is called,
                   then an out of bounds error is returned`,
			rectangle:   domain.NewDrawRectangle(uuid.New(), domain.NewPoint(1, 0), 30, 30, ' ', ' ', time.Now().UTC()),
			expectedErr: domain.ErrOutOfBounds,
		},
		{
			name: `Given a canvas with dimensions 30x30 and a draw rectangle with dimensions 30x30 from point 0x1,
                   when the AddDrawRectangle method is called,
                   then an out of bounds error is returned`,
			rectangle:   domain.NewDrawRectangle(uuid.New(), domain.NewPoint(0, 1), 30, 30, ' ', ' ', time.Now().UTC()),
			expectedErr: domain.ErrOutOfBounds,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			canvas := validCanvas()
			err := canvas.AddDrawRectangle(tt.rectangle)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCanvas_AddFill(t *testing.T) {
	tests := []struct {
		name        string
		fill        domain.Fill
		expectedErr error
	}{
		{
			name: `Given a canvas with dimensions 30x30 and a fill operation from point 0x0,
                   when the AddDrawRectangle method is called,
                   then no error is returned`,
			fill: domain.NewFill(uuid.New(), domain.NewPoint(0, 0), ' ', time.Now().UTC()),
		},
		{
			name: `Given a canvas with dimensions 30x30 and a fill operation from point 30x30,
                   when the AddDrawRectangle method is called,
                   then no error is returned`,
			fill: domain.NewFill(uuid.New(), domain.NewPoint(30, 30), ' ', time.Now().UTC()),
		},
		{
			name: `Given a canvas with dimensions 30x30 and a fill operation from point 31x30,
                   when the AddDrawRectangle method is called,
                   then an out of bounds error is returned`,
			fill:        domain.NewFill(uuid.New(), domain.NewPoint(31, 30), ' ', time.Now().UTC()),
			expectedErr: domain.ErrOutOfBounds,
		},
		{
			name: `Given a canvas with dimensions 30x30 and a fill operation from point 30x31,
                   when the AddDrawRectangle method is called,
                   then an out of bounds error is returned`,
			fill:        domain.NewFill(uuid.New(), domain.NewPoint(30, 31), ' ', time.Now().UTC()),
			expectedErr: domain.ErrOutOfBounds,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			canvas := validCanvas()
			err := canvas.AddFill(tt.fill)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
