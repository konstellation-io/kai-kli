package server

import (
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/spf13/cobra"
)

const (
	_clientSecretFlag = "client-secret"
	_usernameFlag     = "username"
	_passwordFlag     = "password"
)

// NewLoginCLICmd creates a new command to log in to an existing server.
func NewLoginCLICmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use: `login-cli <server_name> --realm <realm> --username <username> --password <password> --client-id <client-id>` +
			`--client-secret <client-secret>`,
		Args:  cobra.ExactArgs(1),
		Short: "login to an existing server",
		Example: `
    $ kli server login-cli my-server --username my-user --password my-password--client-id my-client-id --client-secret my-client-secret
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := args[0]

			realm, err := cmd.Flags().GetString(_realmFlag)
			if err != nil {
				return err
			}
			clientID, err := cmd.Flags().GetString(_clientIDFlag)
			if err != nil {
				return err
			}
			clientSecret, err := cmd.Flags().GetString(_clientSecretFlag)
			if err != nil {
				return err
			}
			username, err := cmd.Flags().GetString(_usernameFlag)
			if err != nil {
				return err
			}
			password, err := cmd.Flags().GetString(_passwordFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			cfg, err := configuration.NewKaiConfigService(logger).GetConfiguration()
			if err != nil {
				return err
			}

			s, err := cfg.GetServer(serverName)
			if err != nil {
				return err
			}

			authenticator := authserver.NewCLIAuthenticator(logger, gocloak.NewClient(s.AuthEndpoint))

			_, err = server.NewHandler(logger, r, authenticator).LoginCLI(serverName, realm, clientID, clientSecret, username, password)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Logged in to %q.", serverName))

			return nil
		},
	}

	cmd.Flags().String(_realmFlag, "", "Realm to login.")
	cmd.Flags().String(_clientIDFlag, "", "Client ID to login.")
	cmd.Flags().String(_clientSecretFlag, "", "Client secret to login.")
	cmd.Flags().String(_usernameFlag, "", "Username to login.")
	cmd.Flags().String(_passwordFlag, "", "Password to login.")

	cmd.MarkFlagRequired(_realmFlag)        //nolint:errcheck
	cmd.MarkFlagRequired(_clientIDFlag)     //nolint:errcheck
	cmd.MarkFlagRequired(_clientSecretFlag) //nolint:errcheck
	cmd.MarkFlagRequired(_usernameFlag)     //nolint:errcheck
	cmd.MarkFlagRequired(_passwordFlag)     //nolint:errcheck

	return cmd
}
