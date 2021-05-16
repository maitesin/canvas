package ascii

import "errors"

var (
	ErrInvalidTask        = errors.New("invalid task")
	ErrRendersOutOfBounds = errors.New("rendering out of bounds")
)
