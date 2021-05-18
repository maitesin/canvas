package ascii_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/domain"
	"github.com/maitesin/sketch/internal/infra/ascii"
	"github.com/stretchr/testify/require"
)

func canvasFixture1() domain.Canvas {
	rectangle1 := domain.NewDrawRectangle(
		uuid.New(),
		domain.NewPoint(3, 2),
		3,
		5,
		'X',
		'@',
		time.Now().UTC(),
	)
	rectangle2 := domain.NewDrawRectangle(
		uuid.New(),
		domain.NewPoint(10, 3),
		6,
		14,
		'0',
		'X',
		time.Now().UTC(),
	)
	return domain.NewCanvas(
		uuid.New(),
		9,
		24,
		[]domain.Task{rectangle1, rectangle2},
		time.Now().UTC(),
	)
}

func outputFixture1() string {
	return `                        
                        
   @@@@@                
   @XXX@  XXXXXXXXXXXXXX
   @@@@@  X000000000000X
          X000000000000X
          X000000000000X
          X000000000000X
          XXXXXXXXXXXXXX
`
}

func canvasFixture2() domain.Canvas {
	rectangle1 := domain.NewDrawRectangle(
		uuid.New(),
		domain.NewPoint(14, 0),
		6,
		7,
		'.',
		'.',
		time.Now().UTC(),
	)
	rectangle2 := domain.NewDrawRectangle(
		uuid.New(),
		domain.NewPoint(0, 3),
		4,
		8,
		' ',
		'O',
		time.Now().UTC(),
	)
	rectangle3 := domain.NewDrawRectangle(
		uuid.New(),
		domain.NewPoint(5, 5),
		3,
		5,
		'X',
		'X',
		time.Now().UTC(),
	)
	return domain.NewCanvas(
		uuid.New(),
		8,
		21,
		[]domain.Task{rectangle1, rectangle2, rectangle3},
		time.Now().UTC(),
	)
}

func outputFixture2() string {
	return `              .......
              .......
              .......
OOOOOOOO      .......
O      O      .......
O    XXXXX    .......
OOOOOXXXXX           
     XXXXX           
`
}

func canvasFixture3(t *testing.T) domain.Canvas {
	t.Helper()

	fill := domain.NewFill(
		uuid.New(),
		domain.NewPoint(0, 0),
		'-',
		time.Now().UTC(),
	)
	canvas := canvasFixture2()
	err := canvas.AddFill(fill)
	require.NoError(t, err)
	return canvas
}

func outputFixture3() string {
	return `--------------.......
--------------.......
--------------.......
OOOOOOOO------.......
O      O------.......
O    XXXXX----.......
OOOOOXXXXX-----------
     XXXXX-----------
`
}

func invalidCanvas() domain.Canvas {
	return domain.NewCanvas(
		uuid.New(),
		10,
		10,
		[]domain.Task{"this will fail"},
		time.Now().UTC(),
	)
}

func TestRenderer(t *testing.T) {
	tests := []struct {
		name           string
		canvas         domain.Canvas
		expectedOutput string
		expectedErr    error
	}{
		{
			name: `Given the canvas from the fixture 1,
                   when the render method is called from the ASCII renderer,
                   then it outputs the output shown in the description of the task`,
			canvas:         canvasFixture1(),
			expectedOutput: outputFixture1(),
		},
		{
			name: `Given the canvas from the fixture 2,
                   when the render method is called from the ASCII renderer,
                   then it outputs the output shown in the description of the task`,
			canvas:         canvasFixture2(),
			expectedOutput: outputFixture2(),
		},
		{
			name: `Given the canvas from the fixture 3,
                   when the render method is called from the ASCII renderer,
                   then it outputs the output shown in the description of the task`,
			canvas:         canvasFixture3(t),
			expectedOutput: outputFixture3(),
		},
		{
			name: `Given a canvas with an invalid task,
                   when the render method is called from the ASCII renderer,
                   then it will return an error of invalid task`,
			canvas:      invalidCanvas(),
			expectedErr: ascii.ErrInvalidTask,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			re := ascii.Renderer{}
			writer := &bytes.Buffer{}
			err := re.Render(writer, tt.canvas)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedOutput, writer.String())
			}
		})
	}
}
