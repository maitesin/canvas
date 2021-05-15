package http

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type CreateCanvasRequest struct {
	ID uuid.UUID `json:"id"`
}

type RequestType string

const (
	DrawRectangleRequestType RequestType = "draw_rectangle"
	AddFillRequestType       RequestType = "add_fill"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type DrawRectangleRequest struct {
	ID      uuid.UUID `json:"id"`
	Point   Point     `json:"point"`
	Height  int       `json:"height"`
	Width   int       `json:"width"`
	Filler  *string   `json:"filler,omitempty"`
	Outline *string   `json:"outline,omitempty"`
}

func (drr DrawRectangleRequest) Validate() error {
	if drr.Filler == nil && drr.Outline == nil {
		return errors.New("both filler and outline cannot be empty. Once of them must be present")
	}
	if drr.Filler != nil && len(*drr.Filler) != 1 {
		return errors.New("filler must be a single character")
	}
	if drr.Outline != nil && len(*drr.Outline) != 1 {
		return errors.New("outline must be a single character")
	}

	return nil
}

type AddFillRequest struct {
	ID     uuid.UUID `json:"id"`
	Point  Point     `json:"point"`
	Filler string    `json:"filler"`
}

func (afr AddFillRequest) Validate() error {
	if len(afr.Filler) != 1 {
		return errors.New("filler must be a single character")
	}

	return nil
}

type TaskRequest struct {
	Type      RequestType           `json:"type"`
	Rectangle *DrawRectangleRequest `json:"rectangle,omitempty"`
	Fill      *AddFillRequest       `json:"fill,omitempty"`
}

func (tr TaskRequest) Validate() error {
	switch tr.Type {
	case DrawRectangleRequestType:
		if tr.Rectangle == nil {
			return fmt.Errorf("rectangle attribute must be present in task %q", DrawRectangleRequestType)
		}
		return tr.Rectangle.Validate()
	case AddFillRequestType:
		if tr.Fill == nil {
			return fmt.Errorf("fill attribute must be present in task %q", AddFillRequestType)
		}
		return tr.Fill.Validate()
	default:
		return fmt.Errorf("unsupported operation %q", tr.Type)
	}
}
