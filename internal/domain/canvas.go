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

// Area defines a square area in a 2D space
type Area struct {
	point  Point
	height int
	width  int
}

// NewArea is a constructor for areas
func NewArea(point Point, height, width int) Area {
	return Area{
		point:  point,
		height: height,
		width:  width,
	}
}

// Rectangle defines the coordinates of a 2D rectangle with all its properties
type Rectangle struct {
	id      uuid.UUID
	area    Area
	filler  rune
	outline rune

	createdAt time.Time
}

// NewRectangle is a constructor for rectangles
func NewRectangle(id uuid.UUID, area Area, filler, outline rune, createdAt time.Time) Rectangle {
	return Rectangle{
		id:        id,
		area:      area,
		filler:    filler,
		outline:   outline,
		createdAt: createdAt,
	}
}

// Canvas defines the 2D space where rectangles can be placed
type Canvas struct {
	id     uuid.UUID
	area   Area
	layers []Rectangle

	createdAt time.Time
}

// NewCanvas is a constructor for canvas
func NewCanvas(id uuid.UUID, area Area, layers []Rectangle, createdAt time.Time) Canvas {
	return Canvas{
		id:        id,
		area:      area,
		layers:    layers,
		createdAt: createdAt,
	}
}

// AddRectangle adds a rectangle to an existing canvas
func (c *Canvas) AddRectangle(rectangle Rectangle) {
	c.layers = append(c.layers, rectangle)
}
