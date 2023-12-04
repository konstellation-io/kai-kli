package server

import (
	"fmt"

	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
)

const _defaultFlag = "default"
const _insecureFlag = "insecure"

// NewAddCmd creates a new command to add a new server to config file.
func NewAddCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <server_name> <server_host> [--default]",
		Aliases: []string{"set"},
		Args:    cobra.ExactArgs(2), //nolint:gomnd
		Short:   "Add a new server in config file",
		Example: `
    $ kli server add my-server local.kai
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			setDefault, err := cmd.Flags().GetBool(_defaultFlag)
			if err != nil {
				return err
			}

			isInsecure, err := cmd.Flags().GetBool(_insecureFlag)
			if err != nil {
				return err
			}

			newServer := &server.Server{
				Name:     args[0],
				Host:     args[1],
				Protocol: getServerProtocol(isInsecure),
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
	cmd.Flags().Bool(_insecureFlag, false, "Set server's protocol to http.")

	return cmd
}

func getServerProtocol(isInsecure bool) string {
	if isInsecure {
		return server.ProtocolHTTP
	}

	return server.ProtocolHTTPS
}
