// Package transport provides HTTP transport utilities for writing JSON responses.
package transport

import(
	"errors"
)


var (
	ErrInvalidRequest        = errors.New("invalid request")
	ErrEmptyReqBody          = errors.New("request body is empty")
	ErrFailedToDecodeReqBody = errors.New("failed to decode request body")
	ErrEncode = errors.New("failed to encode JSON") 
)
