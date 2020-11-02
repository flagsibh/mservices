package data

import (
	"fmt"
)

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// ErrGenericError is a generic error message returned by a server
type ErrGenericError struct {
	Message string `json:"message"`
}

// ErrGenericErrors is a collection of generic or validation errors.
type ErrGenericErrors struct {
	Messages []string `json:"messages"`
}
