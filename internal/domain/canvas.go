package domain

import (
	"time"

	"github.com/google/uuid"
)

// Point defines a point in a 2D space
type Point struct {
	x, y uint
}

// NewPoint is a constructor for points
func NewPoint(x, y uint) Point {
	return Point{
		x: x,
		y: y,
	}
}

// Area defines a square area in a 2D space
type Area struct {
	point  Point
	height uint
	width  uint
}

// NewArea is a constructor for areas
func NewArea(point Point, height, width uint) Area {
	return Area{
		point:  point,
		height: height,
		width:  width,
	}
}

type Task interface{}

// DrawRectangle defines the coordinates of how to draw a 2D rectangle
type DrawRectangle struct {
	id      uuid.UUID
	area    Area
	filler  rune
	outline rune

	createdAt time.Time
}

// NewDrawRectangle is a constructor for tasks that draw rectangles
func NewDrawRectangle(id uuid.UUID, area Area, filler, outline rune, createdAt time.Time) DrawRectangle {
	return DrawRectangle{
		id:        id,
		area:      area,
		filler:    filler,
		outline:   outline,
		createdAt: createdAt,
	}
}

// Fill defines the coordinates where a filling operation needs to be performed
type Fill struct {
	id     uuid.UUID
	point  Point
	filler rune

	createdAt time.Time
}

func NewFill(id uuid.UUID, point Point, filler rune, createdAt time.Time) Fill {
	return Fill{
		id:        id,
		point:     point,
		filler:    filler,
		createdAt: createdAt,
	}
}

// Canvas defines the 2D space where rectangles can be placed
type Canvas struct {
	id     uuid.UUID
	height uint
	width  uint
	tasks  []Task

	createdAt time.Time
}

// Tasks returns the slice of tasks to be performed in the canvas
func (c Canvas) Tasks() []Task {
	return c.tasks
}

// NewCanvas is a constructor for canvas
func NewCanvas(id uuid.UUID, height, width uint, tasks []Task, createdAt time.Time) Canvas {
	return Canvas{
		id:        id,
		height:    height,
		width:     width,
		tasks:     tasks,
		createdAt: createdAt,
	}
}

// AddDrawRectangle adds a rectangle to an existing canvas
func (c *Canvas) AddDrawRectangle(rectangle DrawRectangle) {
	c.tasks = append(c.tasks, rectangle)
}

// AddFill adds a fill operation to an existing canvas
func (c *Canvas) AddFill(fill Fill) {
	c.tasks = append(c.tasks, fill)
}
