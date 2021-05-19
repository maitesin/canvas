package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
	log "github.com/sirupsen/logrus" //nolint: depguard
)

func CreateCanvasHandler(handler app.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := LoggerFromContext(r.Context())
		var createCanvasRequest CreateCanvasRequest
		if err := json.NewDecoder(r.Body).Decode(&createCanvasRequest); err != nil {
			logger.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cmd := app.CreateCanvasCmd{
			ID: createCanvasRequest.ID,
		}

		if err := handler.Handle(r.Context(), cmd); err != nil {
			logger.WithField("canvas_id", createCanvasRequest.ID).Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("http://%s%s/%s", r.Host, r.RequestURI, createCanvasRequest.ID.String()))
		w.WriteHeader(http.StatusCreated)
	}
}

func AddTaskHandler(drawRectangle, addFill app.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := LoggerFromContext(r.Context())
		var taskRequest TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
			logger.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := taskRequest.Validate(); err != nil {
			logger.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		canvasID, err := uuid.Parse(chi.URLParam(r, "canvasID"))
		if err != nil {
			logger.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cmd := createCmdFromTaskRequest(taskRequest, canvasID)

		handler := drawRectangle
		if taskRequest.Type == AddFillRequestType {
			handler = addFill
		}

		if err := handler.Handle(r.Context(), cmd); err != nil {
			logger.WithField("canvas_id", canvasID).Error(err)
			switch {
			case errors.As(err, &app.CanvasNotFound{}):
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			case errors.Is(err, domain.ErrOutOfBounds):
				http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func createCmdFromTaskRequest(request TaskRequest, canvasID uuid.UUID) app.Command {
	switch request.Type {
	case DrawRectangleRequestType:
		return createDrawRectangleCmdFromTaskRequest(request, canvasID)
	case AddFillRequestType:
		return createAddFillCmdFromTaskRequest(request, canvasID)
	}

	return nil
}

func createDrawRectangleCmdFromTaskRequest(request TaskRequest, canvasID uuid.UUID) app.Command {
	filler := ' '
	if request.Rectangle.Filler != nil {
		filler = []rune(*request.Rectangle.Filler)[0]
	}

	outline := filler
	if request.Rectangle.Outline != nil {
		outline = []rune(*request.Rectangle.Outline)[0]
	}

	return app.DrawRectangleCmd{
		CanvasID:    canvasID,
		RectangleID: request.Rectangle.ID,
		Point: domain.NewPoint(
			request.Rectangle.Point.X,
			request.Rectangle.Point.Y,
		),
		Height:  request.Rectangle.Height,
		Width:   request.Rectangle.Width,
		Filler:  filler,
		Outline: outline,
	}
}

func createAddFillCmdFromTaskRequest(request TaskRequest, canvasID uuid.UUID) app.Command {
	return app.AddFillCmd{
		CanvasID: canvasID,
		FillID:   request.Fill.ID,
		Point: domain.NewPoint(
			request.Fill.Point.X,
			request.Fill.Point.Y,
		),
		Filler: []rune(request.Fill.Filler)[0],
	}
}

func RenderCanvasHandler(handler app.QueryHandler, renderer Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := LoggerFromContext(r.Context())

		canvasID, err := uuid.Parse(chi.URLParam(r, "canvasID"))
		if err != nil {
			logger.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		loggerFields := log.Fields{
			"canvas_id": canvasID,
		}

		query := app.RetrieveCanvasQuery{
			ID: canvasID,
		}

		queryResponse, err := handler.Handle(r.Context(), query)
		if err != nil {
			logger.WithFields(loggerFields).Error(err)
			switch {
			case errors.As(err, &app.CanvasNotFound{}):
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		canvas, ok := queryResponse.(domain.Canvas)
		if !ok {
			logger.WithFields(loggerFields).Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = renderer.Render(w, canvas)
		if err != nil {
			logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
