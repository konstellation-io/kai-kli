package args

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/pkg/errors"
)

func TestCheckServerFlagError(t *testing.T) {
	cmd := &cobra.Command{
		Args: CheckServerFlag,
		Run:  func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().String("server", "", "server flag")
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.EqualError(t, err, errors.ErrNoServerConf.Error())

	cmd.SetArgs([]string{"--server"})
	err = cmd.Execute()

	require.EqualError(t, err, fmt.Errorf("flag needs an argument: --server").Error())
}

func TestCheckServerFlag(t *testing.T) {
	testArgs := []string{"--server", "test"}

	cmd := &cobra.Command{}
	cmd.Flags().String("server", "", "server flag")
	cmd.SetArgs(testArgs)

	err := cmd.Execute()
	require.NoError(t, err)
}

func TestComposeArgsCheck(t *testing.T) {
	composed := ComposeArgsCheck(cobra.ExactArgs(1), CheckServerFlag)
	cmd := &cobra.Command{
		Args: composed,
		Run:  func(_ *cobra.Command, _ []string) {},
	}
	cmd.SetArgs([]string{"some", "args"})
	err := cmd.Execute()
	require.EqualError(t, err, fmt.Errorf("accepts 1 arg(s), received 2").Error())
}
