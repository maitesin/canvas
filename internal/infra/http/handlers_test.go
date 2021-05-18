package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	log "github.com/sirupsen/logrus" //nolint: depguard
	"github.com/stretchr/testify/require"
)

func validCommandHandler() app.CommandHandler {
	return &CommandHandlerMock{
		HandleFunc: func(context.Context, app.Command) error {
			return nil
		},
	}
}

type commandHandlerMutator func(app.CommandHandler) app.CommandHandler

var noopCommandHandlerMutator = func(commandHandler app.CommandHandler) app.CommandHandler { return commandHandler }

func validCreateCanvasBodyReader(t *testing.T, id uuid.UUID) io.Reader {
	t.Helper()

	request := httpx.CreateCanvasRequest{
		ID: id,
	}

	b, err := json.Marshal(request)
	require.NoError(t, err)

	return bytes.NewReader(b)
}

func TestCreateCanvasHandler(t *testing.T) {
	createdCanvasID := uuid.New()

	tests := []struct {
		name                        string
		commandHandlerMutator       commandHandlerMutator
		bodyReader                  io.Reader
		expectedStatusCode          int
		expectedLocationHeaderValue string
	}{
		{
			name: `Given a working command handler, a valid canvas ID, and a valid body request,
                   when the add fill handler is called,
                   then a status created (201) response is returned`,
			commandHandlerMutator:       noopCommandHandlerMutator,
			bodyReader:                  validCreateCanvasBodyReader(t, createdCanvasID),
			expectedStatusCode:          http.StatusCreated,
			expectedLocationHeaderValue: fmt.Sprintf("http:///%s", createdCanvasID.String()),
		},
		{
			name: `Given a working command handler, a valid canvas ID, and an invalid body request,
                   when the add fill handler is called,
                   then a status bad request (400) response is returned`,
			commandHandlerMutator: noopCommandHandlerMutator,
			bodyReader:            func() io.Reader { return strings.NewReader("") }(),
			expectedStatusCode:    http.StatusBadRequest,
		},
		{
			name: `Given a non-working command handler, a valid canvas ID, and a valid body request,
                   when the add fill handler is called,
                   then a status internal server error (500) response is returned`,
			commandHandlerMutator: func(app.CommandHandler) app.CommandHandler {
				handler := &CommandHandlerMock{
					HandleFunc: func(context.Context, app.Command) error {
						return errors.New("something else went wrong")
					},
				}
				return handler
			},
			bodyReader:         validCreateCanvasBodyReader(t, uuid.New()),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			commandHandler := tt.commandHandlerMutator(validCommandHandler())

			logger := log.New()
			ctx := httpx.ContextWithLogger(context.Background(), logger)

			req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/canvas/", tt.bodyReader)
			require.NoError(t, err)

			res := httptest.NewRecorder()

			httpx.CreateCanvasHandler(commandHandler)(res, req)
			result := res.Result()
			defer result.Body.Close()

			require.Equal(t, tt.expectedStatusCode, result.StatusCode)
			require.Equal(t, tt.expectedLocationHeaderValue, result.Header.Get("Location"))
		})
	}
}

func validAddFillerBodyReader(t *testing.T) io.Reader {
	t.Helper()

	b, err := json.Marshal(validAddFillRequest())
	require.NoError(t, err)

	return bytes.NewReader(b)
}

func validDrawRectangleBodyReader(t *testing.T) io.Reader {
	t.Helper()

	b, err := json.Marshal(validDrawRectangleRequest())
	require.NoError(t, err)

	return bytes.NewReader(b)
}

func invalidDrawRectangleBodyReader(t *testing.T) io.Reader {
	t.Helper()

	b, err := json.Marshal(invalidDrawRectangleRequestMissingBothFillerAndOutline())
	require.NoError(t, err)

	return bytes.NewReader(b)
}

// nolint:funlen
func TestAddTaskHandler(t *testing.T) {
	tests := []struct {
		name                  string
		commandHandlerMutator commandHandlerMutator
		canvasID              string
		bodyReader            io.Reader
		expectedStatusCode    int
	}{
		{
			name: `Given a working command handler, a valid canvas ID, and a valid body request,
                   when the add task handler is called,
                   then a status ok (200) response is returned`,
			commandHandlerMutator: noopCommandHandlerMutator,
			canvasID:              uuid.New().String(),
			bodyReader:            validAddFillerBodyReader(t),
			expectedStatusCode:    http.StatusOK,
		},
		{
			name: `Given a working command handler, a valid canvas ID, and another valid body request,
                   when the add task handler is called,
                   then a status ok (200) response is returned`,
			commandHandlerMutator: noopCommandHandlerMutator,
			canvasID:              uuid.New().String(),
			bodyReader:            validDrawRectangleBodyReader(t),
			expectedStatusCode:    http.StatusOK,
		},
		{
			name: `Given a working command handler, an invalid canvas ID, and a valid body request,
                   when the add task handler is called,
                   then a status bad request (400) response is returned`,
			commandHandlerMutator: noopCommandHandlerMutator,
			canvasID:              "wololo",
			bodyReader:            validAddFillerBodyReader(t),
			expectedStatusCode:    http.StatusBadRequest,
		},
		{
			name: `Given a working command handler, a valid canvas ID, and a non-JSON body request,
                   when the add task handler is called,
                   then a status bad request (400) response is returned`,
			commandHandlerMutator: noopCommandHandlerMutator,
			canvasID:              uuid.New().String(),
			bodyReader:            func() io.Reader { return strings.NewReader("") }(),
			expectedStatusCode:    http.StatusBadRequest,
		},
		{
			name: `Given a working command handler, a valid canvas ID, and an invalid body request,
                   when the add task handler is called,
                   then a status bad request (400) response is returned`,
			commandHandlerMutator: noopCommandHandlerMutator,
			canvasID:              uuid.New().String(),
			bodyReader:            invalidDrawRectangleBodyReader(t),
			expectedStatusCode:    http.StatusBadRequest,
		},
		{
			name: `Given a working command handler, a valid canvas ID (but not associated with any canvas), and a valid body request,
                   when the add task handler is called,
                   then a status not found (404) response is returned`,
			commandHandlerMutator: func(app.CommandHandler) app.CommandHandler {
				handler := &CommandHandlerMock{
					HandleFunc: func(context.Context, app.Command) error {
						return app.CanvasNotFound{}
					},
				}
				return handler
			},
			canvasID:           uuid.New().String(),
			bodyReader:         validDrawRectangleBodyReader(t),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: `Given a non-working command handler, a valid canvas ID, and a valid body request,
                   when the add task handler is called,
                   then a status internal server error (500) response is returned`,
			commandHandlerMutator: func(app.CommandHandler) app.CommandHandler {
				handler := &CommandHandlerMock{
					HandleFunc: func(context.Context, app.Command) error {
						return errors.New("something else went wrong")
					},
				}
				return handler
			},
			canvasID:           uuid.New().String(),
			bodyReader:         validDrawRectangleBodyReader(t),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			commandHandler := tt.commandHandlerMutator(validCommandHandler())

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("canvasID", tt.canvasID)
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)
			logger := log.New()
			ctx = httpx.ContextWithLogger(ctx, logger)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("/canvas/%s/fill", tt.canvasID), tt.bodyReader)
			require.NoError(t, err)

			res := httptest.NewRecorder()

			httpx.AddTaskHandler(commandHandler, noopCommandHandlerMutator(validCommandHandler()))(res, req)
			result := res.Result()
			defer result.Body.Close()

			require.Equal(t, tt.expectedStatusCode, result.StatusCode)
		})
	}
}

func validQueryHandler() app.QueryHandler {
	return &QueryHandlerMock{
		HandleFunc: func(context.Context, app.Query) (app.QueryResponse, error) {
			return domain.Canvas{}, nil
		},
	}
}

type queryHandlerMutator func(app.QueryHandler) app.QueryHandler

var noopQueryHandlerMutator = func(handler app.QueryHandler) app.QueryHandler { return handler }

//nolint:funlen
func TestRenderCanvasHandler(t *testing.T) {
	tests := []struct {
		name                string
		queryHandlerMutator queryHandlerMutator
		rendererMutator     rendererMutator
		canvasID            string
		expectedStatusCode  int
	}{
		{
			name: `Given a working query handler, a working renderer, and a valid canvas ID,
                   when the render canvas handler is called,
                   then a status code ok (200) is returned`,
			queryHandlerMutator: noopQueryHandlerMutator,
			rendererMutator:     noopRendererMutator,
			canvasID:            uuid.New().String(),
			expectedStatusCode:  http.StatusOK,
		},
		{
			name: `Given a working query handler, a working renderer, and an invalid canvas ID,
                   when the render canvas handler is called,
                   then a status code bad request (400) is returned`,
			queryHandlerMutator: noopQueryHandlerMutator,
			rendererMutator:     noopRendererMutator,
			canvasID:            "",
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name: `Given a working query handler, a working renderer, and a valid canvas ID, but that does not exists,
                   when the render canvas handler is called,
                   then a status code not found (404) is returned`,
			queryHandlerMutator: func(app.QueryHandler) app.QueryHandler {
				handler := &QueryHandlerMock{
					HandleFunc: func(context.Context, app.Query) (app.QueryResponse, error) {
						return nil, app.CanvasNotFound{}
					},
				}
				return handler
			},
			rendererMutator:    noopRendererMutator,
			canvasID:           uuid.New().String(),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: `Given a non-working query handler, a working renderer, and a valid canvas ID,
                   when the render canvas handler is called,
                   then a status code internal server error (500) is returned`,
			queryHandlerMutator: func(app.QueryHandler) app.QueryHandler {
				handler := &QueryHandlerMock{
					HandleFunc: func(context.Context, app.Query) (app.QueryResponse, error) {
						return nil, errors.New("something else went wrong")
					},
				}
				return handler
			},
			rendererMutator:    noopRendererMutator,
			canvasID:           uuid.New().String(),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: `Given a working query handler, but that returns a string, a working renderer, and a valid canvas ID,
                   when the render canvas handler is called,
                   then a status code internal server error (500) is returned`,
			queryHandlerMutator: func(app.QueryHandler) app.QueryHandler {
				handler := &QueryHandlerMock{
					HandleFunc: func(context.Context, app.Query) (app.QueryResponse, error) {
						return "", nil
					},
				}
				return handler
			},
			rendererMutator:    noopRendererMutator,
			canvasID:           uuid.New().String(),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: `Given a working query handler, a non-working renderer, and a valid canvas ID,
                   when the render canvas handler is called,
                   then a status code internal server error (500) is returned`,
			queryHandlerMutator: noopQueryHandlerMutator,
			rendererMutator: func(httpx.Renderer) httpx.Renderer {
				renderer := &RendererMock{
					RenderFunc: func(io.Writer, domain.Canvas) error {
						return errors.New("something went wrong")
					},
				}
				return renderer
			},
			canvasID:           uuid.New().String(),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			queryHandler := tt.queryHandlerMutator(validQueryHandler())
			renderer := tt.rendererMutator(validRenderer())

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("canvasID", tt.canvasID)
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)
			logger := log.New()
			ctx = httpx.ContextWithLogger(ctx, logger)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/canvas/%s", uuid.New()), nil)
			require.NoError(t, err)

			res := httptest.NewRecorder()

			httpx.RenderCanvasHandler(queryHandler, renderer)(res, req)
			result := res.Result()
			defer result.Body.Close()

			require.Equal(t, tt.expectedStatusCode, result.StatusCode)
		})
	}
}
