package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/domain"
)

//go:generate moq -out zmock_command_test.go -pkg app_test . Command

// Command defines the interface of the commands to be performed
type Command interface {
	Name() string
}

//go:generate moq -out ../infra/http/zmock_command_test.go -pkg http_test . CommandHandler

// CommandHandler defines the interface of the handler to run commands
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CreateCanvasCmd is a VTO
type CreateCanvasCmd struct {
	ID uuid.UUID
}

// Name returns the name of the command to create a canvas
func (c CreateCanvasCmd) Name() string {
	return "createCanvas"
}

// CreateCanvasHandler is the handler to create a canvas
type CreateCanvasHandler struct {
	repository   CanvasRepository
	canvasHeight uint
	canvasWidth  uint
}

// NewCreateCanvasHandler is a constructor
func NewCreateCanvasHandler(repository CanvasRepository) CreateCanvasHandler {
	return CreateCanvasHandler{repository: repository}
}

// Handle creates a canvas
func (c CreateCanvasHandler) Handle(ctx context.Context, cmd Command) error {
	createCmd, ok := cmd.(CreateCanvasCmd)
	if !ok {
		return InvalidCommandError{expected: CreateCanvasCmd{}, received: cmd}
	}

	canvas := domain.NewCanvas(
		createCmd.ID,
		c.canvasHeight,
		c.canvasWidth,
		nil,
		time.Now().UTC(),
	)

	return c.repository.Insert(ctx, canvas)
}

// DrawRectangleCmd is a VTO
type DrawRectangleCmd struct {
	CanvasID    uuid.UUID
	RectangleID uuid.UUID
	Point       domain.Point
	Height      uint
	Width       uint
	Filler      rune
	Outline     rune
}

// Name returns the name of the command to draw a rectangle in a canvas
func (c DrawRectangleCmd) Name() string {
	return "drawRectangle"
}

// DrawRectangleHandler is the handler to draw a rectangle in a canvas
type DrawRectangleHandler struct {
	repository CanvasRepository
}

// NewDrawRectangleHandler is a constructor
func NewDrawRectangleHandler(repository CanvasRepository) DrawRectangleHandler {
	return DrawRectangleHandler{repository: repository}
}

// Handle adds a draw rectangle task to a canvas
func (d DrawRectangleHandler) Handle(ctx context.Context, cmd Command) error {
	drawRectangleCmd, ok := cmd.(DrawRectangleCmd)
	if !ok {
		return InvalidCommandError{expected: DrawRectangleCmd{}, received: cmd}
	}

	canvas, err := d.repository.FindByID(ctx, drawRectangleCmd.CanvasID)
	if err != nil {
		return err
	}

	rectangle := domain.NewDrawRectangle(
		drawRectangleCmd.RectangleID,
		domain.NewArea(drawRectangleCmd.Point, drawRectangleCmd.Height, drawRectangleCmd.Width),
		drawRectangleCmd.Filler,
		drawRectangleCmd.Outline,
		time.Now().UTC(),
	)

	canvas.AddDrawRectangle(rectangle)

	return d.repository.Update(ctx, canvas)
}

// AddFillCmd is a VTO
type AddFillCmd struct {
	CanvasID uuid.UUID
	FillID   uuid.UUID
	Point    domain.Point
	Filler   rune
}

// Name returns the name of the command to add a fill task in a canvas
func (c AddFillCmd) Name() string {
	return "fill"
}

// AddFillHandler is the handler to add a fill in a canvas
type AddFillHandler struct {
	repository CanvasRepository
}

// NewAddFillHandler is a constructor
func NewAddFillHandler(repository CanvasRepository) AddFillHandler {
	return AddFillHandler{repository: repository}
}

// Handle adds a fill task to a canvas
func (f AddFillHandler) Handle(ctx context.Context, cmd Command) error {
	addFillCmd, ok := cmd.(AddFillCmd)
	if !ok {
		return InvalidCommandError{expected: AddFillCmd{}, received: cmd}
	}

	canvas, err := f.repository.FindByID(ctx, addFillCmd.CanvasID)
	if err != nil {
		return err
	}

	fill := domain.NewFill(
		addFillCmd.FillID,
		addFillCmd.Point,
		addFillCmd.Filler,
		time.Now().UTC(),
	)

	canvas.AddFill(fill)

	return f.repository.Update(ctx, canvas)
}
