package errors

import (
	"fmt"
)

func NewSrcPathValidationError(workflow, node, src string) *ValidationError {
	message := fmt.Sprintf("src file %q for node %q in workflow %q does not exists ", src, node, workflow)
	return NewValidationError("nodes.src", message)
}

type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field %q: %s", e.Field, e.Message)
}
