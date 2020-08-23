package errtype

import "errors"

var (
	ErrEmptyContent = errors.New("empty content")
	EmptySegments   = errors.New("empty segments")
)
