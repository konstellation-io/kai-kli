package storage

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_serverFlag = "server"
)

// NewStorageCmd creates a new command to handle 'storage' subcommands.
func NewStorageCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage <command> [opts...]",
		Short: "Manage storage for kli",
		Example: heredoc.Doc(`
			$ kli storage open
		`),
	}

	cmd.PersistentFlags().StringP(_serverFlag, "s", "", "KAI server to use")

	cmd.AddCommand(
		NewOpenCmd(logger),
	)

	return cmd
}
