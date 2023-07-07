package server

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewServerCmd creates a new command to handle 'server' subcommands.
func NewServerCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command>",
		Short: "Manage servers for kli",
		Example: heredoc.Doc(`
			$ kli server ls
			$ kli server add my-server http://api.kai.local
		`),
	}

	cmd.AddCommand(
		NewListCmd(logger),
		NewDefaultCmd(logger, cfg),
		NewAddCmd(logger, cfg),
	)

	return cmd
}
