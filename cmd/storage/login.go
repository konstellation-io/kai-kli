package storage

import (
	"github.com/konstellation-io/kli/internal/commands/storage"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewLoginCmd creates a new command to log in to an existing server.
func NewLoginCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login <server_name> --auth-url <auth_url> --realm <realm> --client-id <client-id>",
		Args:  cobra.ExactArgs(0),
		Short: "login to an existing server",
		Example: `
    $ kli server login my-server --user my-user --password my-password
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				return err
			}

			err = storage.NewHandler(logger).Login(serverName)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
