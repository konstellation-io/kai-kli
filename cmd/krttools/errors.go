package krttools

import "fmt"

func NewValidationError(err error) error {
	return fmt.Errorf("error validating krt: %w", err)
}

func NewBuildError(err error) error {
	return fmt.Errorf("error building krt: %w", err)
}
