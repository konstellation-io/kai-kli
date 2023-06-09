package config

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewConfigCmd manage config command for version.
func NewConfigCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Get or set config values",
	}

	cmd.AddCommand(
		NewListConfigCmd(logger, cfg),
		NewUpdateConfigCmd(logger, cfg),
	)

	return cmd
}
