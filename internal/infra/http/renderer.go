package http

import (
	"io"

	"github.com/maitesin/sketch/internal/domain"
)

//go:generate moq -out zmock_renderer_test.go -pkg http_test . Renderer

type Renderer interface {
	Render(writer io.Writer, canvas domain.Canvas) error
}
