package app

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	errMsgInvalidCommand = "invalid command %q received. Expected %q"
	errMsgInvalidQuery   = "invalid query %q received. Expected %q"
	errMsgCanvasNotFound = "canvas %q not found"
)

type InvalidCommandError struct {
	Received Command
	Expected Command
}

func (ice InvalidCommandError) Error() string {
	return fmt.Sprintf(errMsgInvalidCommand, ice.Received.Name(), ice.Expected.Name())
}

type InvalidQueryError struct {
	Received Query
	Expected Query
}

func (iqe InvalidQueryError) Error() string {
	return fmt.Sprintf(errMsgInvalidQuery, iqe.Received.Name(), iqe.Expected.Name())
}

type CanvasNotFound struct {
	ID uuid.UUID
}

func (cnf CanvasNotFound) Error() string {
	return fmt.Sprintf(errMsgCanvasNotFound, cnf.ID)
}
