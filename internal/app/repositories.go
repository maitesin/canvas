package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/domain"
)

//go:generate moq -out zmock_canvas_repository_test.go -pkg app_test . CanvasRepository
//go:generate moq -out ../infra/http/zmock_canvas_repository_test.go -pkg http_test . CanvasRepository

type CanvasRepository interface {
	Insert(ctx context.Context, canvas domain.Canvas) error
	Update(ctx context.Context, canvas domain.Canvas) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Canvas, error)
}
