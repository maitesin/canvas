package mem

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
)

type CanvasRepository struct {
	m        sync.RWMutex
	canvases []domain.Canvas
}

func (c *CanvasRepository) Insert(_ context.Context, canvas domain.Canvas) error {
	c.m.Lock()
	defer c.m.Unlock()

	c.canvases = append(c.canvases, canvas)

	return nil
}

func (c *CanvasRepository) Update(_ context.Context, canvas domain.Canvas) error {
	c.m.Lock()
	defer c.m.Unlock()

	for i := range c.canvases {
		if c.canvases[i].ID() == canvas.ID() {
			c.canvases[i] = canvas
			return nil
		}
	}

	return app.CanvasNotFound{
		ID: canvas.ID(),
	}
}

func (c *CanvasRepository) FindByID(_ context.Context, id uuid.UUID) (domain.Canvas, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	for i := range c.canvases {
		if c.canvases[i].ID() == id {
			return c.canvases[i], nil
		}
	}

	return domain.Canvas{}, app.CanvasNotFound{ID: id}
}
