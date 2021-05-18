package http_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maitesin/sketch/config"
	"github.com/maitesin/sketch/internal/app"
	"github.com/maitesin/sketch/internal/domain"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	log "github.com/sirupsen/logrus" //nolint: depguard
	"github.com/stretchr/testify/require"
)

func validCanvasRepository() app.CanvasRepository {
	return &CanvasRepositoryMock{
		FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
			return domain.NewCanvas(uuid.New(), 30, 30, nil, time.Now().UTC()), nil
		},
		InsertFunc: func(context.Context, domain.Canvas) error {
			return nil
		},
		UpdateFunc: func(context.Context, domain.Canvas) error {
			return nil
		},
	}
}

type repositoryMutator func(app.CanvasRepository) app.CanvasRepository

var noopRepositoryMutator = func(repository app.CanvasRepository) app.CanvasRepository { return repository }

func validRenderer() httpx.Renderer {
	return &RendererMock{
		RenderFunc: func(io.Writer, domain.Canvas) error {
			return nil
		},
	}
}

type rendererMutator func(httpx.Renderer) httpx.Renderer

var noopRendererMutator = func(renderer httpx.Renderer) httpx.Renderer { return renderer }

func TestDefaultRouter_CreateCanvas(t *testing.T) {
	tests := []struct {
		name               string
		repositoryMutator  repositoryMutator
		bodyReader         io.Reader
		expectedStatusCode int
	}{
		{
			name: `Given a working canvas repository and a valid body request,
                   when the endpoint to create a new canvas is called,
                   then a canvas is successfully created and a status code created (201) is returned`,
			repositoryMutator:  noopRepositoryMutator,
			bodyReader:         strings.NewReader(fmt.Sprintf(`{"id":%q}`, uuid.New().String())),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: `Given a working canvas repository and an invalid body request,
                   when the endpoint to create a new canvas is called,
                   then a status code bad request (400) is returned`,
			repositoryMutator:  noopRepositoryMutator,
			bodyReader:         strings.NewReader(""),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: `Given a non-working canvas repository and a valid body request,
                   when the endpoint to create a new canvas is called,
                   then a status code internal server error (500) is returned`,
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					InsertFunc: func(context.Context, domain.Canvas) error {
						return errors.New("something went wrong")
					},
				}
				return repository
			},
			bodyReader:         strings.NewReader(fmt.Sprintf(`{"id":%q}`, uuid.New().String())),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := tt.repositoryMutator(validCanvasRepository())

			cfg, err := config.New()
			require.NoError(t, err)

			logger := log.New()
			ctx := httpx.ContextWithLogger(context.Background(), logger)

			server := httptest.NewServer(httpx.DefaultRouter(ctx, cfg.Canvas, repository, &RendererMock{}))
			defer server.Close()

			client := server.Client()
			//nolint: noctx
			resp, err := client.Post(fmt.Sprintf("%s/canvas", server.URL), "application/json", tt.bodyReader)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tt.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestDefaultRouter_AddTaskToCanvas(t *testing.T) {
	tests := []struct {
		name               string
		repositoryMutator  repositoryMutator
		bodyReader         io.Reader
		expectedStatusCode int
	}{
		{
			name: `Given a working canvas repository and a valid body request,
                   when the endpoint to add a draw rectangle task to a canvas is called,
                   then a status code OK (200) is returned`,
			repositoryMutator:  noopRepositoryMutator,
			bodyReader:         validDrawRectangleBodyReader(t),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: `Given a working canvas repository and a valid body request,
                   when the endpoint to add a fill task to a canvas is called,
                   then a status code OK (200) is returned`,
			repositoryMutator:  noopRepositoryMutator,
			bodyReader:         validAddFillerBodyReader(t),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: `Given a working canvas repository and an invalid body request,
                   when the endpoint to add a fill task to a canvas is called,
                   then a status code bad request (400) is returned`,
			repositoryMutator:  noopRepositoryMutator,
			bodyReader:         strings.NewReader(""),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: `Given a working canvas repository and a valid body request, but the canvasID present in the URL does not exists,
                   when the endpoint to add a draw rectangle task to a canvas is called,
                   then a status code not found (404) is returned`,
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, app.CanvasNotFound{}
					},
				}
				return repository
			},
			bodyReader:         validDrawRectangleBodyReader(t),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: `Given a non-working canvas repository and a valid body request,
                   when the endpoint to add a draw rectangle task to a canvas is called,
                   then a status code internal server error (500) is returned`,
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, errors.New("something else went wrong")
					},
				}
				return repository
			},
			bodyReader:         validDrawRectangleBodyReader(t),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := tt.repositoryMutator(validCanvasRepository())

			cfg, err := config.New()
			require.NoError(t, err)

			logger := log.New()
			ctx := httpx.ContextWithLogger(context.Background(), logger)

			server := httptest.NewServer(httpx.DefaultRouter(ctx, cfg.Canvas, repository, &RendererMock{}))
			defer server.Close()

			client := server.Client()
			//nolint: noctx
			resp, err := client.Post(fmt.Sprintf("%s/canvas/%s", server.URL, uuid.New()), "application/json", tt.bodyReader)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tt.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestDefaultRouter_RenderCanvas(t *testing.T) {
	tests := []struct {
		name               string
		repositoryMutator  repositoryMutator
		rendererMutator    rendererMutator
		expectedStatusCode int
	}{
		{
			name: `Given a working canvas repository and a working renderer,
                   when the endpoint to render a canvas is called,
                   then a status code OK (200) is returned`,
			repositoryMutator:  noopRepositoryMutator,
			rendererMutator:    noopRendererMutator,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: `Given a working canvas repository and a working renderer, but the canvasID does not exists,
                   when the endpoint to render a canvas is called,
                   then a status code not found (404) is returned`,
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, app.CanvasNotFound{}
					},
				}
				return repository
			},
			rendererMutator:    noopRendererMutator,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: `Given a non-working canvas repository and a working renderer,
                   when the endpoint to render a canvas is called,
                   then a status code internal server error (500) is returned`,
			repositoryMutator: func(app.CanvasRepository) app.CanvasRepository {
				repository := &CanvasRepositoryMock{
					FindByIDFunc: func(context.Context, uuid.UUID) (domain.Canvas, error) {
						return domain.Canvas{}, errors.New("something went wrong")
					},
				}
				return repository
			},
			rendererMutator:    noopRendererMutator,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: `Given a working canvas repository and a non-working renderer,
                   when the endpoint to render a canvas is called,
                   then a status code internal server error (500) is returned`,
			repositoryMutator: noopRepositoryMutator,
			rendererMutator: func(httpx.Renderer) httpx.Renderer {
				renderer := &RendererMock{
					RenderFunc: func(io.Writer, domain.Canvas) error {
						return errors.New("unable to render")
					},
				}
				return renderer
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := tt.repositoryMutator(validCanvasRepository())
			renderer := tt.rendererMutator(validRenderer())

			cfg, err := config.New()
			require.NoError(t, err)

			logger := log.New()
			ctx := httpx.ContextWithLogger(context.Background(), logger)

			server := httptest.NewServer(httpx.DefaultRouter(ctx, cfg.Canvas, repository, renderer))
			defer server.Close()

			client := server.Client()
			//nolint: noctx
			resp, err := client.Get(fmt.Sprintf("%s/canvas/%s", server.URL, uuid.New()))
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tt.expectedStatusCode, resp.StatusCode)
		})
	}
}
