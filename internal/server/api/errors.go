package api

import "errors"

var (
	ErrFileNameIsEmpty   = errors.New("file name is empty")
	ErrFileNameIsInvalid = errors.New("file name is invalid")
)
