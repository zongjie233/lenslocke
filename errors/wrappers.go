package errors

import "errors"

// This function takes an error and a boolean value and returns a boolean value
// It checks if the error is an instance of the Error type and if the boolean value is true
// If both conditions are true, it returns true, otherwise it returns false

var (
	As = errors.As
	Is = errors.Is
)
