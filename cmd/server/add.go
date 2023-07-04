package server

import (
	"fmt"

	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/server"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/internal/logging"
)

const _defaultFlag = "default"

// NewAddCmd creates a new command to add a new server to config file.
func NewAddCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <remote_name> <server_url> [--default]",
		Aliases: []string{"set"},
		Args:    cobra.ExactArgs(2), //nolint:gomnd
		Short:   "Add a new server in config file",
		Example: `
    $ kli server add my-server http://api.local.kre
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

			err = server.NewKaiConfigurator(logger, r).AddNewServer(newServer, setDefault)
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
