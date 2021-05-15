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
	"github.com/maitesin/sketch/internal/infra/ascii"
)

func CreateCanvasHandler(handler app.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createCanvasRequest CreateCanvasRequest
		if err := json.NewDecoder(r.Body).Decode(&createCanvasRequest); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cmd := app.CreateCanvasCmd{
			ID: createCanvasRequest.ID,
		}

		if err := handler.Handle(r.Context(), cmd); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("http://%s%s/%s", r.Host, r.RequestURI, createCanvasRequest.ID.String()))
		w.WriteHeader(http.StatusCreated)
	}
}

func AddTaskHandler(drawRectangle, addFill app.CommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskRequest TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := taskRequest.Validate(); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		canvasID, err := uuid.Parse(chi.URLParam(r, "canvasID"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cmd := createCmdFromTaskRequest(taskRequest, canvasID)

		handler := drawRectangle
		if taskRequest.Type == AddFillRequestType {
			handler = addFill
		}

		if err := handler.Handle(r.Context(), cmd); err != nil {
			switch {
			case errors.Is(err, app.CanvasNotFound{}):
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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

func RenderCanvasHandler(handler app.QueryHandler, renderer ascii.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		canvasID, err := uuid.Parse(chi.URLParam(r, "canvasID"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		query := app.RetrieveCanvasQuery{
			ID: canvasID,
		}

		queryResponse, err := handler.Handle(r.Context(), query)
		if err != nil {
			switch {
			case errors.Is(err, app.CanvasNotFound{}):
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		canvas, ok := queryResponse.(domain.Canvas)
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		matrix, err := renderer.Render(canvas)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		for _, line := range matrix {
			fmt.Fprintln(w, line)
		}

		//httpResponse := RetrieveCanvas{
		//	ID:     canvas.ID(),
		//	Height: canvas.Height(),
		//	Width:  canvas.Width(),
		//	Tasks:  len(canvas.Tasks()),
		//}
		//
		//if err := json.NewEncoder(w).Encode(httpResponse); err != nil {
		//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		//}
	}
}
