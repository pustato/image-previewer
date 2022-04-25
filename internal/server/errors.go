package server

import "errors"

var (
	ErrMalformedRequestPath = errors.New("malformed request path")
	ErrWidthIsNotANumber    = errors.New("width is not a number")
	ErrHeightIsNotANumber   = errors.New("height is not a number")
)
