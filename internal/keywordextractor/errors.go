package keywordextractor

import (
	"errors"
)

var (
	ExternalErr      error = errors.New("something went wrong on the external server")
	ErrInvalidOutput error = errors.New("something went wrong with receiving output")
)
