package args

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/pkg/errors"
)

// ComposeArgsCheck gets multiple checker functions and composes a new check.
func ComposeArgsCheck(fn ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, f := range fn {
			err := f(cmd, args)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// CheckServerFlag checks that server is configured by default or passed as a flag.
func CheckServerFlag(cmd *cobra.Command, _ []string) error {
	s, _ := cmd.Flags().GetString("server")

	if s == "" {
		return errors.ErrNoServerConf
	}

	return nil
}
