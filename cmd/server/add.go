package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const _defaultFlag = "default"

// NewAddCmd creates a new command to add a new server to config file.
func NewAddCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <server_name> <server_url> [--default]",
		Aliases: []string{"set"},
		Args:    cobra.ExactArgs(2), //nolint:gomnd
		Short:   "Add a new server in config file",
		Example: `
    $ kli server add my-server http://api.local.kai
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			setDefault, err := cmd.Flags().GetBool(_defaultFlag)
			if err != nil {
				return err
			}

			newServer := server.Server{
				Name: args[0],
				URL:  args[1],
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = server.NewHandler(logger, r).AddNewServer(newServer, setDefault)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Server %q added.", newServer.Name))

			return nil
		},
	}

	cmd.Flags().BoolP(_defaultFlag, "d", false, "Set the current server as the default server.")

	return cmd
}
