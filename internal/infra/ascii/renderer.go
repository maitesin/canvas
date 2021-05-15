package ascii

import (
	"github.com/maitesin/sketch/internal/domain"
)

type Renderer struct{}

func (Renderer) Render(c domain.Canvas) ([]string, error) {
	if c.Height() == 0 {
		return nil, RendersOutOfBoundsErr
	}
	canvas := make([][]rune, c.Height())
	for i := range canvas {
		canvas[i] = make([]rune, c.Width())
		for j := range canvas[i] {
			canvas[i][j] = ' '
		}
	}

	tasks := c.Tasks()
	for i := range tasks {
		rectangle, ok := tasks[i].(domain.DrawRectangle)
		if ok {
			err := drawRectangle(canvas, rectangle)
			if err != nil {
				return nil, err
			}
		} else {
			fill, ok := tasks[i].(domain.Fill)
			if ok {
				err := addFill(canvas, fill)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, InvalidTaskErr
			}
		}
	}

	result := make([]string, len(canvas))
	for i := range canvas {
		result[i] = string(canvas[i])
	}
	return result, nil
}

func drawRectangle(canvas [][]rune, rectangle domain.DrawRectangle) error {
	h := len(canvas)
	y := rectangle.Point().Y()
	if h < y+rectangle.Height() {
		return RendersOutOfBoundsErr
	}

	w := len(canvas[0])
	x := rectangle.Point().X()
	if w < x+rectangle.Width() {
		return RendersOutOfBoundsErr
	}

	for i := x; i < x+rectangle.Width(); i++ {
		canvas[y][i] = rectangle.Outline()
		canvas[y+rectangle.Height()-1][i] = rectangle.Outline()
	}
	for i := y; i < y+rectangle.Height(); i++ {
		canvas[i][x] = rectangle.Outline()
		canvas[i][x+rectangle.Width()-1] = rectangle.Outline()
	}

	for i := x + 1; i < x+rectangle.Width()-1; i++ {
		for j := y + 1; j < y+rectangle.Height()-1; j++ {
			canvas[j][i] = rectangle.Filler()
		}
	}

	return nil
}

func addFill(canvas [][]rune, fill domain.Fill) error {
	return nil
}
