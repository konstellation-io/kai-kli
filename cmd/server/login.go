package server

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	_authURLFlag  = "auth-url"
	_userFlag     = "user"
	_passwordFlag = "password"
	_realmFlag    = "realm"
	_clientIDFlag = "client-id"
)

// NewLoginCmd creates a new command to log in to an existing server.
func NewLoginCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login <server_name> --auth-url <auth_url> --realm <realm> --client-id <client-id>",
		Args:  cobra.ExactArgs(1),
		Short: "login to an existing server",
		Example: `
    $ kli server login my-server --user my-user --password my-password
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := args[0]
			authURL, err := cmd.Flags().GetString(_authURLFlag)

			if strings.HasPrefix(authURL, "http") || strings.HasPrefix(authURL, "https") {
				return fmt.Errorf("invalid authentication URL, the URL must have the protocol (http:// or https://)")
			}

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

			_, err = server.NewHandler(logger, r).Login(serverName, authURL, realm, clientID)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Logged in to %q.", serverName))

			return nil
		},
	}

	cmd.Flags().String(_authURLFlag, "", "URL to login.")
	cmd.Flags().String(_realmFlag, "", "Realm to login.")
	cmd.Flags().String(_clientIDFlag, "", "ClientID to login.")

	cmd.MarkFlagRequired(_authURLFlag)  //nolint:errcheck
	cmd.MarkFlagRequired(_realmFlag)    //nolint:errcheck
	cmd.MarkFlagRequired(_clientIDFlag) //nolint:errcheck

	return cmd
}
