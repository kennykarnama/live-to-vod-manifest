package errtype

import "errors"

var (
	ErrNotValidScheme = errors.New("invalid scheme url")
	ErrNotS3          = errors.New("URL not a S3 storage")
)
