package storage

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewStorageCmd creates a new command to handle 'storage' subcommands.
func NewStorageCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage <command> [opts...]",
		Short: "Manage storage for kli",
		Example: heredoc.Doc(`
			$ kli server login <server_name> <username> <password>
		`),
	}

	cmd.AddCommand(
		NewLoginCmd(logger),
	)

	return cmd
}
