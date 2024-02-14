package server

import (
	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

const (
	_realmFlag    = "realm"
	_clientIDFlag = "client-id"
)

// NewLoginCmd creates a new command to log in to an existing server.
func NewLoginCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login --realm <realm> --client-id <client-id> [--server <server_name>]",
		Args:  cobra.ExactArgs(0),
		Short: "login to an existing server",
		Example: `
    $ kli server login --realm <realm> --client-id <client-id> --server my-server
		`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			srv, err := getServerOrDefault(cmd, logger)
			if err != nil {
				return err
			}

			realm, err := cmd.Flags().GetString(_realmFlag)
			if err != nil {
				return err
			}
			clientID, err := cmd.Flags().GetString(_clientIDFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			_, err = server.NewHandler(logger, r, authserver.NewDefaultAuthServer(logger)).Login(srv.Name, realm, clientID)
			if err != nil {
				return err
			}

			r.RenderLogin(srv.Name)

			return nil
		},
	}

	cmd.Flags().String(_realmFlag, "", "Realm to login.")
	cmd.Flags().String(_clientIDFlag, "", "ClientID to login.")
	cmd.Flags().StringP(_serverFlag, "s", "", "server name, defaults to default server configuration")

	cmd.MarkFlagRequired(_realmFlag)    //nolint:errcheck
	cmd.MarkFlagRequired(_clientIDFlag) //nolint:errcheck

	return cmd
}
