package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	server2 "github.com/konstellation-io/kli/internal/server"
)

// NewAddCmd creates a new command to remove an existing server from the config file.
func NewRemoveCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <remote_name>",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1), //nolint:gomnd
		Short:   "Removes an existing server from the config file",
		Example: `
    $ kli server rm my-server
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			server := args[0]

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err := server2.NewServerHandler(logger, r).RemoveServer(server)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Server %q removed.", server))

			return nil
		},
	}

	cmd.Flags().BoolP(_defaultFlag, "d", false, "Set the current server as the default server.")

	return cmd
}
