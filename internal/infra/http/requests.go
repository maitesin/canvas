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
	X uint `json:"x"`
	Y uint `json:"y"`
}

type DrawRectangleRequest struct {
	ID      uuid.UUID `json:"id"`
	Point   Point     `json:"point"`
	Height  uint      `json:"height"`
	Width   uint      `json:"width"`
	Filler  *rune     `json:"filler,omitempty"`
	Outline *rune     `json:"outline,omitempty"`
}

func (drr DrawRectangleRequest) Validate() error {
	if drr.Filler == nil && drr.Outline == nil {
		return errors.New("both filler and outline cannot be empty. Once of them must be present")
	}
	return nil
}

type AddFillRequest struct {
	ID     uuid.UUID `json:"id"`
	Point  Point     `json:"point"`
	Filler rune      `json:"filler"`
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
	default:
		return fmt.Errorf("unsupported operation %q", tr.Type)
	}

	return nil
}
