package ascii

import "errors"

var (
	InvalidTaskErr        = errors.New("invalid task")
	RendersOutOfBoundsErr = errors.New("rendering out of bounds")
)
