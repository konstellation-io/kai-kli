package server

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewServerCmd creates a new command to handle 'server' subcommands.
func NewServerCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command> [opts...]",
		Short: "Manage servers for kli",
		Example: heredoc.Doc(`
			$ kli server ls
			$ kli server add <server_name> <server_url>
			$ kli server remove <server_name>
			$ kli server default <server_name>
			$ kli server login <server_name> <username> <password>
			$ kli server logout <server_name>
		`),
	}

	cmd.AddCommand(
		NewListCmd(logger),
		NewDefaultCmd(logger),
		NewAddCmd(logger),
		NewRemoveCmd(logger),
		NewLoginCmd(logger),
		NewLogoutCmd(logger),
	)

	return cmd
}
