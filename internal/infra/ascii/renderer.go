package ascii

import (
	"fmt"
	"io"

	"github.com/maitesin/sketch/internal/domain"
)

type Renderer struct{}

func (Renderer) Render(writer io.Writer, c domain.Canvas) error {
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
			drawRectangle(canvas, rectangle)
		} else {
			fill, ok := tasks[i].(domain.Fill)
			if ok {
				addFill(canvas, fill)
			} else {
				return ErrInvalidTask
			}
		}
	}

	for i := range canvas {
		_, err := fmt.Fprintln(writer, string(canvas[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

func drawRectangle(canvas [][]rune, rectangle domain.DrawRectangle) {
	y := rectangle.Point().Y()
	x := rectangle.Point().X()

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
}

func addFill(canvas [][]rune, fill domain.Fill) {
	flood(canvas, fill.Point(), canvas[fill.Point().Y()][fill.Point().X()], fill.Filler())
}

func flood(canvas [][]rune, point domain.Point, old, new rune) {
	if len(canvas) <= point.Y() || len(canvas[0]) <= point.X() || point.Y() < 0 || point.X() < 0 {
		return
	}
	if canvas[point.Y()][point.X()] == old {
		canvas[point.Y()][point.X()] = new
		flood(canvas, domain.NewPoint(point.X()+1, point.Y()), old, new)
		flood(canvas, domain.NewPoint(point.X()-1, point.Y()), old, new)
		flood(canvas, domain.NewPoint(point.X(), point.Y()+1), old, new)
		flood(canvas, domain.NewPoint(point.X(), point.Y()-1), old, new)
	}
}
