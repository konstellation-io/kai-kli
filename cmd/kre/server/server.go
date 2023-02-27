package server

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewServerCmd creates a new command to handle 'server' subcommands.
func NewServerCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command>",
		Short: "Manage servers for kli",
		Example: heredoc.Doc(`
			$ kli server ls
			$ kli server add my-server http://api.kre.local TOKEN_12345
			$ kli server default my-server
		`),
	}

	cmd.AddCommand(
		NewListCmd(logger, cfg),
		NewDefaultCmd(logger, cfg),
		NewAddCmd(logger, cfg),
	)

	return cmd
}
