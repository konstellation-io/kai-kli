package server

import (
	"fmt"

	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewGetTokenCmd creates a new command to retrieve server session token.
func NewGetTokenCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "get-token <server_name>",
		Args:   cobra.ExactArgs(1),
		Short:  "Retrieve server session token.",
		Hidden: true,
		Example: `
    $ kli server get-token my-server
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			token, err := server.NewHandler(logger, r, authserver.NewDefaultAuthServer(logger)).GetToken(args[0])
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Token for server '%s' retrieved successfully: %s", args[0], token.AccessToken))

			return nil
		},
	}

	return cmd
}
