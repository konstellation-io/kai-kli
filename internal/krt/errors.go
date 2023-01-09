package krt

import (
	"fmt"
	"strings"
)

//// ValidationErrors is an array of Error's to store multiple error messages post validation.
//// NOTE: this technique based on github.com/go-playground/validator.
// type ValidationErrors []error
//
//// Error method allows ValidationErrors to subscribe to the Error interface.
// func (ve ValidationErrors) Error() string {
//	buff := bytes.NewBufferString("")
//
//	var e error
//
//	for i := 0; i < len(ve); i++ {
//		e = ve[i]
//		buff.WriteString(e.Error())
//		buff.WriteString("\n")
//	}
//
//	return strings.TrimSpace(buff.String())
//}

func newSrcPathValidationError(workflow, node, src string) string {
	return fmt.Sprintf("error src file %q for node %q in workflow %q does not exists ", src, node, workflow)
}

type ValidationError struct {
	Messages []string
}

func NewValidationError(errs []error) ValidationError {
	return ValidationError{
		Messages: getErrorMessages(errs),
	}
}

func getErrorMessages(errs []error) []string {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	return messages
}

func (e ValidationError) Error() string {
	return strings.Join(e.Messages, "\n")
}
