package http

import "github.com/google/uuid"

type RetrieveCanvas struct {
	ID     uuid.UUID `json:"id"`
	Height uint      `json:"height"`
	Width  uint      `json:"width"`
	Tasks  uint      `json:"tasks"`
}
