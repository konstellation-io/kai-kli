package storage

import (
	"github.com/konstellation-io/kli/internal/commands/storage"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewOpenCmd creates a new command to log in to an existing server.
func NewOpenCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Args:  cobra.ExactArgs(0),
		Short: "open storage console webpage for given sever",
		Example: `
    $ kli storage open
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				return err
			}

			err = storage.NewHandler(logger).OpenConsole(serverName)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
