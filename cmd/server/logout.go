package server

import (
	"errors"

	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
)

var (
	ErrServerNotFound = errors.New("server not found")
)

// NewLogoutCmd creates a new command to log out from an existing server.
func NewLogoutCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout [--server <server_name>]",
		Args:  cobra.ExactArgs(0),
		Short: "logout from an existing server",
		Example: `
    $ kli server logout --server my-server
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := getServerOrDefault(cmd, logger)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = server.NewHandler(logger, r, authserver.NewDefaultAuthServer(logger)).Logout(srv.Name)
			if err != nil {
				return err
			}

			r.RenderLogout(srv.Name)

			return nil
		},
	}

	cmd.Flags().StringP(_serverFlag, "s", "", "server name, defaults to default server configuration")

	return cmd
}
