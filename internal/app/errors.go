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
	received Command
	expected Command
}

func (ice InvalidCommandError) Error() string {
	return fmt.Sprintf(errMsgInvalidCommand, ice.received.Name(), ice.expected.Name())
}

type InvalidQueryError struct {
	received Query
	expected Query
}

func (iqe InvalidQueryError) Error() string {
	return fmt.Sprintf(errMsgInvalidQuery, iqe.received.Name(), iqe.expected.Name())
}

type CanvasNotFound struct {
	ID uuid.UUID
}

func (cnf CanvasNotFound) Error() string {
	return fmt.Sprintf(errMsgCanvasNotFound, cnf.ID)
}
