package errtype

import "errors"

var (
	ErrEmptyCsvRecord = errors.New("empty csv record")
)
