package errors

import "errors"

// ErrNotInitialized is returned when cobra.Command is not initialized
var ErrNotInitialized = errors.New("not initialized")

// ErrTargetNotPointer godoc
var ErrTargetNotPointer = errors.New("target is not pointer")

// ErrFileDoesNotExist godoc
var ErrFileDoesNotExist = errors.New("file does not exist")
