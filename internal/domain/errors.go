package domain

import "errors"

// ErrOutOfBounds used when a draw rectangle or a fill operation does not fit into the canvas
var ErrOutOfBounds = errors.New("task out of bounds")
