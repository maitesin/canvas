package http

import "github.com/google/uuid"

type RetrieveCanvas struct {
	ID     uuid.UUID `json:"id"`
	Height int       `json:"height"`
	Width  int       `json:"width"`
	Tasks  int       `json:"tasks"`
}
