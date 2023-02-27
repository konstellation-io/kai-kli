package testhelpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// SetupConfigDir setup an XDG_CONFIG_HOME for testing.
func SetupConfigDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "kli-test")
	require.NoError(t, err)

	err = os.Setenv("XDG_CONFIG_HOME", dir)
	require.NoError(t, err)

	return dir
}

// CleanConfigDir clean temporal testing XDG_CONFIG_HOME.
func CleanConfigDir(t *testing.T, dir string) {
	t.Helper()

	err := os.RemoveAll(dir)
	require.NoError(t, err)
}
