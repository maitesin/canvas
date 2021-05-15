package app

import (
	"context"

	"github.com/google/uuid"
)

//go:generate moq -out zmock_query_test.go -pkg app_test . Query

// Query defines the interface of the queries to be performed
type Query interface {
	Name() string
}

// QueryResponse defines the response to be received from the QueryHandler
type QueryResponse interface{}

// QueryHandler defines the interface of the handler to run queries
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (QueryResponse, error)
}

// RetrieveCanvasQuery is a VTO
type RetrieveCanvasQuery struct {
	ID uuid.UUID
}

// Name returns the name of the query to retrieve a canvas
func (r RetrieveCanvasQuery) Name() string {
	return "retrieveCanvas"
}

// RetrieveCanvasHandler is the handler to retrieve the canvas
type RetrieveCanvasHandler struct {
	repository CanvasRepository
}

// NewRetrieveCanvasHandler is a constructor
func NewRetrieveCanvasHandler(repository CanvasRepository) RetrieveCanvasHandler {
	return RetrieveCanvasHandler{repository: repository}
}

// Handle retrieves the canvas
func (r RetrieveCanvasHandler) Handle(ctx context.Context, query Query) (QueryResponse, error) {
	retrieveQuery, ok := query.(RetrieveCanvasQuery)
	if !ok {
		return nil, InvalidQueryError{expected: RetrieveCanvasQuery{}, received: query}
	}

	return r.repository.FindByID(ctx, retrieveQuery.ID)
}
