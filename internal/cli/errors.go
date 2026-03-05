package cli

import (
	"fmt"
	"os"
)

const (
	ExitSuccess    = 0
	ExitError      = 1
	ExitValidation = 2
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

func ExitWithError(err error) {
	if err == nil {
		return
	}

	if validationErr, ok := err.(*ValidationError); ok {
		fmt.Fprintf(os.Stderr, "validation error: %s\n", validationErr.Error())
		os.Exit(ExitValidation)
	}

	fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	os.Exit(ExitError)
}

func ExitWithSuccess(message string) {
	if message != "" {
		fmt.Println(message)
	}
	os.Exit(ExitSuccess)
}

func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
