package http_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	httpx "github.com/maitesin/sketch/internal/infra/http"
	"github.com/stretchr/testify/require"
)

func validDrawRectangleRequest() httpx.TaskRequest {
	filler := "X"
	outline := "O"
	return httpx.TaskRequest{
		Type: httpx.DrawRectangleRequestType,
		Rectangle: &httpx.DrawRectangleRequest{
			ID: uuid.New(),
			Point: httpx.Point{
				X: 10,
				Y: 10,
			},
			Height:  5,
			Width:   5,
			Filler:  &filler,
			Outline: &outline,
		},
	}
}

func invalidDrawRectangleRequestMissingBothFillerAndOutline() httpx.TaskRequest {
	return httpx.TaskRequest{
		Type: httpx.DrawRectangleRequestType,
		Rectangle: &httpx.DrawRectangleRequest{
			ID: uuid.New(),
			Point: httpx.Point{
				X: 10,
				Y: 10,
			},
			Height: 5,
			Width:  5,
		},
	}
}

func invalidDrawRectangleRequestWithFillerTooLong() httpx.TaskRequest {
	filler := "wololo"
	outline := "O"
	return httpx.TaskRequest{
		Type: httpx.DrawRectangleRequestType,
		Rectangle: &httpx.DrawRectangleRequest{
			ID: uuid.New(),
			Point: httpx.Point{
				X: 10,
				Y: 10,
			},
			Height:  5,
			Width:   5,
			Filler:  &filler,
			Outline: &outline,
		},
	}
}

func invalidDrawRectangleRequestWithOutlineTooLong() httpx.TaskRequest {
	filler := "X"
	outline := "wololo"
	return httpx.TaskRequest{
		Type: httpx.DrawRectangleRequestType,
		Rectangle: &httpx.DrawRectangleRequest{
			ID: uuid.New(),
			Point: httpx.Point{
				X: 10,
				Y: 10,
			},
			Height:  5,
			Width:   5,
			Filler:  &filler,
			Outline: &outline,
		},
	}
}

func validAddFillRequest() httpx.TaskRequest {
	return httpx.TaskRequest{
		Type: httpx.AddFillRequestType,
		Fill: &httpx.AddFillRequest{
			ID: uuid.New(),
			Point: httpx.Point{
				X: 5,
				Y: 5,
			},
			Filler: "-",
		},
	}
}

func invalidAddFillRequest() httpx.TaskRequest {
	return httpx.TaskRequest{
		Type: httpx.AddFillRequestType,
		Fill: &httpx.AddFillRequest{
			ID: uuid.New(),
			Point: httpx.Point{
				X: 5,
				Y: 5,
			},
			Filler: "wololo",
		},
	}
}

func TestTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		taskRequest httpx.TaskRequest
		expectedErr error
	}{
		{
			name: `Given a valid draw rectangle request,
                   when the validate method is called,
                   then no error is returned`,
			taskRequest: validDrawRectangleRequest(),
		},
		{
			name: `Given a valid add fill request,
                   when the validate method is called,
                   then no error is returned`,
			taskRequest: validAddFillRequest(),
		},
		{
			name: `Given an invalid draw rectangle request because it has both filler and outline missing,
                   when the validate method is called,
                   then an error is returned`,
			taskRequest: invalidDrawRectangleRequestMissingBothFillerAndOutline(),
			expectedErr: errors.New(""),
		},
		{
			name: `Given an invalid draw rectangle request because it has the filler with too many runes,
                   when the validate method is called,
                   then an error is returned`,
			taskRequest: invalidDrawRectangleRequestWithFillerTooLong(),
			expectedErr: errors.New(""),
		},
		{
			name: `Given an invalid draw rectangle request because it has the outline with too many runes,
                   when the validate method is called,
                   then an error is returned`,
			taskRequest: invalidDrawRectangleRequestWithOutlineTooLong(),
			expectedErr: errors.New(""),
		},
		{
			name: `Given an invalid add fill request because it has the filler with too many runes,
                   when the validate method is called,
                   then an error is returned`,
			taskRequest: invalidAddFillRequest(),
			expectedErr: errors.New(""),
		},
		{
			name: `Given an invalid task request because the task type is invalid,
                   when the validate method is called,
                   then an error is returned`,
			taskRequest: httpx.TaskRequest{},
			expectedErr: errors.New(""),
		},
		{
			name: `Given an invalid task request because the task type is draw rectangle, but the rectangle attribute is missing,
                   when the validate method is called,
                   then an error is returned`,
			taskRequest: httpx.TaskRequest{
				Type: httpx.DrawRectangleRequestType,
			},
			expectedErr: errors.New(""),
		},
		{
			name: `Given an invalid task request because the task type is add fill, but the fill attribute is missing,
                   when the validate method is called,
                   then an error is returned`,
			taskRequest: httpx.TaskRequest{
				Type: httpx.AddFillRequestType,
			},
			expectedErr: errors.New(""),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.taskRequest.Validate()
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
