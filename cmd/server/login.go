package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	_userFlag     = "user"
	_passwordFlag = "password"
	_realmFlag    = "realm"
	_clientIDFlag = "client-id"
)

// NewLoginCmd creates a new command to log in to an existing server.
func NewLoginCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use: "login <remote_name> [--user <username>] " +
			"[--password <password>] [--realm <realm>] [--client-id <client-id>]",
		Args:  cobra.ExactArgs(1),
		Short: "auth.go to an existing server",
		Example: `
    $ kli server auth.go my-server --user my-user --password my-password
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := args[0]
			username, err := cmd.Flags().GetString(_userFlag)
			if err != nil {
				return err
			}
			password, err := cmd.Flags().GetString(_passwordFlag)
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

			_, err = server.NewServerHandler(logger, r).Login(serverName, realm, clientID, username, password)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Logged in to %q.", serverName))

			return nil
		},
	}

	cmd.Flags().StringP(_userFlag, "u", "", "Username to auth.go.")
	cmd.Flags().StringP(_passwordFlag, "p", "", "Password to auth.go.")
	cmd.Flags().String(_realmFlag, "", "Password to auth.go.")
	cmd.Flags().String(_clientIDFlag, "", "Password to auth.go.")

	cmd.MarkFlagRequired(_userFlag)     //nolint:errcheck
	cmd.MarkFlagRequired(_passwordFlag) //nolint:errcheck
	cmd.MarkFlagRequired(_realmFlag)    //nolint:errcheck
	cmd.MarkFlagRequired(_clientIDFlag) //nolint:errcheck

	return cmd
}
