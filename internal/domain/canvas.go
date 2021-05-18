package domain

import (
	"time"

	"github.com/google/uuid"
)

// Point defines a point in a 2D space
type Point struct {
	x, y int
}

// NewPoint is a constructor for points
func NewPoint(x, y int) Point {
	return Point{
		x: x,
		y: y,
	}
}

// X returns the X position of the point
func (p Point) X() int {
	return p.x
}

// Y returns the Y position of the point
func (p Point) Y() int {
	return p.y
}

type Task interface{}

// DrawRectangle defines the coordinates of how to draw a 2D rectangle
type DrawRectangle struct {
	id      uuid.UUID
	point   Point
	height  int
	width   int
	filler  rune
	outline rune

	createdAt time.Time
}

// ID returns the id of the rectangle
func (dr DrawRectangle) ID() uuid.UUID {
	return dr.id
}

// Point returns the point where the rectangle starts
func (dr DrawRectangle) Point() Point {
	return dr.point
}

// Height returns the height of the rectangle
func (dr DrawRectangle) Height() int {
	return dr.height
}

// Width returns the width of the rectangle
func (dr DrawRectangle) Width() int {
	return dr.width
}

// Filler returns the rune to fill the rectangle with
func (dr DrawRectangle) Filler() rune {
	return dr.filler
}

// Outline returns the rune to out line the rectangle with
func (dr DrawRectangle) Outline() rune {
	return dr.outline
}

// CreatedAt returns the time where the rectangle was created
func (dr DrawRectangle) CreatedAt() time.Time {
	return dr.createdAt
}

// NewDrawRectangle is a constructor for tasks that draw rectangles
func NewDrawRectangle(id uuid.UUID, point Point, height, width int, filler, outline rune, createdAt time.Time) DrawRectangle {
	return DrawRectangle{
		id:        id,
		point:     point,
		height:    height,
		width:     width,
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

// ID returns the id of the fill
func (f Fill) ID() uuid.UUID {
	return f.id
}

// Point returns the point where the fill starts
func (f Fill) Point() Point {
	return f.point
}

// Filler returns the rune to fill the area with
func (f Fill) Filler() rune {
	return f.filler
}

// CreatedAt returns the time where the fill was created
func (f Fill) CreatedAt() time.Time {
	return f.createdAt
}

// NewFill is a constructor
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
	height int
	width  int
	tasks  []Task

	createdAt time.Time
}

// ID returns the ID of the canvas
func (c Canvas) ID() uuid.UUID {
	return c.id
}

// Height returns the height of the canvas
func (c Canvas) Height() int {
	return c.height
}

// Width returns the width of the canvas
func (c Canvas) Width() int {
	return c.width
}

// Tasks returns the slice of tasks to be performed in the canvas
func (c Canvas) Tasks() []Task {
	return c.tasks
}

// CreatedAt returns the time where the canvas was created
func (c Canvas) CreatedAt() time.Time {
	return c.createdAt
}

// NewCanvas is a constructor for canvas
func NewCanvas(id uuid.UUID, height, width int, tasks []Task, createdAt time.Time) Canvas {
	return Canvas{
		id:        id,
		height:    height,
		width:     width,
		tasks:     tasks,
		createdAt: createdAt,
	}
}

// AddDrawRectangle adds a rectangle to an existing canvas
func (c *Canvas) AddDrawRectangle(rectangle DrawRectangle) error {
	if c.height < rectangle.height+rectangle.point.y ||
		c.width < rectangle.width+rectangle.point.x {
		return ErrOutOfBounds
	}

	c.tasks = append(c.tasks, rectangle)
	return nil
}

// AddFill adds a fill operation to an existing canvas
func (c *Canvas) AddFill(fill Fill) error {
	if c.height < fill.point.y || c.width < fill.point.x {
		return ErrOutOfBounds
	}

	c.tasks = append(c.tasks, fill)
	return nil
}
