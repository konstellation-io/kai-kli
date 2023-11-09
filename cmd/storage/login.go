package storage

import (
	"fmt"
	"strings"

	"github.com/konstellation-io/kli/internal/commands/storage"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_authURLFlag = "auth-url"
)

// NewLoginCmd creates a new command to log in to an existing server.
func NewLoginCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login <server_name> --auth-url <auth_url> --realm <realm> --client-id <client-id>",
		Args:  cobra.ExactArgs(0),
		Short: "login to an existing server",
		Example: `
    $ kli server login my-server --user my-user --password my-password
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			authURL, err := cmd.Flags().GetString(_authURLFlag)
			if err != nil {
				return err
			}

			if !strings.HasPrefix(authURL, "http://") && !strings.HasPrefix(authURL, "https://") {
				return fmt.Errorf("invalid authentication URL, the URL must have the protocol (http:// or https://)")
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = storage.NewHandler(logger, r).Login(authURL)
			if err != nil {
				return err
			}

			logger.Success("Logged in!")

			return nil
		},
	}

	cmd.Flags().String(_authURLFlag, "", "URL to login.")

	cmd.MarkFlagRequired(_authURLFlag) //nolint:errcheck

	return cmd
}
