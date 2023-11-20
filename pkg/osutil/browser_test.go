package osutil_test

import (
	"errors"
	"os/exec"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/konstellation-io/kli/pkg/osutil"
	"github.com/stretchr/testify/assert"
)

var errRunCmd = errors.New("error cmd run")

func mockExecCommand(wantRunError bool) func(string, ...string) *exec.Cmd {
	return func(_ string, _ ...string) *exec.Cmd {
		cmd := &exec.Cmd{}

		_ = monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(c *exec.Cmd) error {
			if wantRunError {
				return errRunCmd
			}

			return nil
		})

		return cmd
	}
}

func TestOpenBrowser(t *testing.T) {
	_ = monkey.Patch(exec.Command, mockExecCommand(false))
	defer monkey.UnpatchAll()

	err := osutil.OpenBrowser("testing.com")
	assert.NoError(t, err)
}

func TestOpenBrowser_ErrorExecutingCommand(t *testing.T) {
	_ = monkey.Patch(exec.Command, mockExecCommand(true))
	defer monkey.UnpatchAll()

	err := osutil.OpenBrowser("testing.com")
	assert.ErrorIs(t, err, errRunCmd)
}
