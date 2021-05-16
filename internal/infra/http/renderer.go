package http

import (
	"io"

	"github.com/maitesin/sketch/internal/domain"
)

type Renderer interface {
	Render(writer io.Writer, canvas domain.Canvas) error
}
