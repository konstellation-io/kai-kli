package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/server"
)

// NewRemoveCmd creates a new command to remove an existing server from the config file.
func NewRemoveCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <server_name>",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Short:   "Removes an existing server from the config file",
		Example: `
    $ kli server rm my-server
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := args[0]

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err := server.NewServerHandler(logger, r).RemoveServer(serverName)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Server %q removed.", serverName))

			return nil
		},
	}

	return cmd
}
