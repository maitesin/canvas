package sql

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
	"github.com/upper/db/v4"
)

const (
	canvasTable     = "canvases"
	rectanglesTable = "rectangles"
	fillsTable      = "fills"
)

func onConflictDoNothing(queryIn string) string {
	return queryIn + `ON CONFLICT DO NOTHING`
}

type Canvas struct {
	ID        uuid.UUID `db:"id"`
	Height    int       `db:"height"`
	Width     int       `db:"width"`
	CreatedAt time.Time `db:"created_at"`
}

type Rectangle struct {
	ID        uuid.UUID `db:"id"`
	CanvasID  uuid.UUID `db:"canvas_id"`
	X         int       `db:"x"`
	Y         int       `db:"y"`
	Height    int       `db:"height"`
	Width     int       `db:"width"`
	Filler    rune      `db:"filler"`
	Outline   rune      `db:"outline"`
	CreatedAt time.Time `db:"created_at"`
}

type Fill struct {
	ID        uuid.UUID `db:"id"`
	CanvasID  uuid.UUID `db:"canvas_id"`
	X         int       `db:"x"`
	Y         int       `db:"y"`
	Filler    rune      `db:"filler"`
	CreatedAt time.Time `db:"created_at"`
}

type CanvasRepository struct {
	sess db.Session
}

func NewCanvasRepository(sess db.Session) *CanvasRepository {
	return &CanvasRepository{sess: sess}
}

func (c *CanvasRepository) Insert(ctx context.Context, canvas domain.Canvas) error {
	sqlCanvas, _, _, err := domainToSQL(canvas)
	if err != nil {
		return err
	}

	_, err = c.sess.WithContext(ctx).
		SQL().
		InsertInto(canvasTable).
		Values(sqlCanvas).
		Amend(onConflictDoNothing).
		Exec()
	return err
}

func domainToSQL(canvas domain.Canvas) (Canvas, []Rectangle, []Fill, error) {
	var rectangles []Rectangle
	var fills []Fill

	for _, task := range canvas.Tasks() {
		rectangle, ok := task.(domain.DrawRectangle)
		if !ok {
			fill, ok := task.(domain.Fill)
			if !ok {
				return Canvas{}, nil, nil, fmt.Errorf("failed to convert domain to sql task: %#v", task)
			}
			fills = append(fills, Fill{
				ID:        fill.ID(),
				CanvasID:  canvas.ID(),
				X:         fill.Point().X(),
				Y:         fill.Point().Y(),
				Filler:    fill.Filler(),
				CreatedAt: fill.CreatedAt(),
			})
			continue
		}
		rectangles = append(rectangles, Rectangle{
			ID:        rectangle.ID(),
			CanvasID:  canvas.ID(),
			X:         rectangle.Point().X(),
			Y:         rectangle.Point().Y(),
			Height:    rectangle.Height(),
			Width:     rectangle.Width(),
			Filler:    rectangle.Filler(),
			Outline:   rectangle.Outline(),
			CreatedAt: rectangle.CreatedAt(),
		})
	}

	return Canvas{
		ID:        canvas.ID(),
		Height:    canvas.Height(),
		Width:     canvas.Width(),
		CreatedAt: canvas.CreatedAt(),
	}, rectangles, fills, nil
}

func (c *CanvasRepository) Update(ctx context.Context, canvas domain.Canvas) error {
	_, sqlRectangles, sqlFills, err := domainToSQL(canvas)
	if err != nil {
		return err
	}

	return c.sess.Tx(func(sess db.Session) error {
		if len(sqlRectangles) > 0 {
			rectanglesInserter := c.sess.WithContext(ctx).
				SQL().
				InsertInto(rectanglesTable)

			for i := range sqlRectangles {
				rectanglesInserter = rectanglesInserter.Values(sqlRectangles[i])
			}
			_, err = rectanglesInserter.
				Amend(onConflictDoNothing).
				Exec()
			if err != nil {
				return err
			}
		}

		if len(sqlFills) > 0 {
			fillsInserter := c.sess.WithContext(ctx).
				SQL().
				InsertInto(fillsTable)

			for i := range sqlFills {
				fillsInserter = fillsInserter.Values(sqlFills[i])
			}

			_, err = fillsInserter.
				Amend(onConflictDoNothing).
				Exec()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *CanvasRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Canvas, error) {
	var sqlCanvas Canvas
	var sqlRectangles []Rectangle
	var sqlFills []Fill

	err := c.sess.Tx(func(sess db.Session) error {
		err := c.sess.WithContext(ctx).
			Collection(canvasTable).
			Find(db.Cond{"id": id}).
			One(&sqlCanvas)
		if err != nil {
			if err == db.ErrNoMoreRows {
				return app.CanvasNotFound{ID: id}
			}
			return err
		}

		err = c.sess.WithContext(ctx).
			Collection(rectanglesTable).
			Find(db.Cond{"canvas_id": id}).
			All(&sqlRectangles)
		if err != nil && err != db.ErrNoMoreRows {
			return err
		}

		err = c.sess.WithContext(ctx).
			Collection(fillsTable).
			Find(db.Cond{"canvas_id": id}).
			All(&sqlFills)
		if err != nil && err != db.ErrNoMoreRows {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Canvas{}, err
	}

	return sqlToDomain(sqlCanvas, sqlRectangles, sqlFills), nil
}

func sqlToDomain(canvas Canvas, rectangles []Rectangle, fills []Fill) domain.Canvas {
	sort.Slice(rectangles, func(i, j int) bool {
		return rectangles[i].CreatedAt.Before(rectangles[j].CreatedAt)
	})
	sort.Slice(fills, func(i, j int) bool {
		return fills[i].CreatedAt.Before(fills[j].CreatedAt)
	})

	tasks := make([]domain.Task, len(rectangles)+len(fills))
	ti := 0
	ri := 0
	fi := 0

	for ri < len(rectangles) && fi < len(fills) {
		if rectangles[ri].CreatedAt.Before(fills[fi].CreatedAt) {
			tasks[ti] = sqlRectangleToDomain(rectangles[ri])
			ri++
		} else {
			tasks[ti] = sqlFillToDomain(fills[fi])
			fi++
		}
		ti++
	}
	for ri < len(rectangles) {
		tasks[ti] = sqlRectangleToDomain(rectangles[ri])
		ri++
		ti++
	}
	for fi < len(fills) {
		tasks[ti] = sqlFillToDomain(fills[fi])
		fi++
		ti++
	}

	return domain.NewCanvas(
		canvas.ID,
		canvas.Height,
		canvas.Width,
		tasks,
		canvas.CreatedAt,
	)
}

func sqlRectangleToDomain(rectangle Rectangle) domain.DrawRectangle {
	return domain.NewDrawRectangle(
		rectangle.ID,
		domain.NewPoint(rectangle.X, rectangle.Y),
		rectangle.Height,
		rectangle.Width,
		rectangle.Filler,
		rectangle.Outline,
		rectangle.CreatedAt,
	)
}

func sqlFillToDomain(fill Fill) domain.Fill {
	return domain.NewFill(
		fill.ID,
		domain.NewPoint(fill.X, fill.Y),
		fill.Filler,
		fill.CreatedAt,
	)
}
