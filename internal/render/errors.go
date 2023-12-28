package render

import (
	"errors"
	"fmt"
)

var (
	// ErrCreatingLogsFile used when server name is unknown.
	ErrCreatingLogsFile = errors.New("error creating logs.txt file")
)

func ErrorCreatingLogsFile(err error) string {
	return fmt.Sprintf("%s: %s", ErrCreatingLogsFile.Error(), err.Error())
}
