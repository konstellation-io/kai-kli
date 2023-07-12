package server

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewServerCmd creates a new command to handle 'server' subcommands.
func NewServerCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command>",
		Short: "Manage servers for kli",
		Example: heredoc.Doc(`
			$ kli server ls
			$ kli server add my-server http://api.kai.local
			$ kli server remove my-server
			$ kli server default my-server
		`),
	}

	cmd.AddCommand(
		NewListCmd(logger),
		NewDefaultCmd(logger),
		NewAddCmd(logger),
		NewRemoveCmd(logger),
	)

	return cmd
}
