package server

import (
	"fmt"

	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewDefaultCmd creates a new command to handle 'default' keyword.
func NewDefaultCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default <server-name>",
		Short: "Set a default server",
		Long:  "Mark an existing server as default",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err := server.NewHandler(logger, r, authserver.NewDefaultAuthServer(logger)).SetDefaultServer(name)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Server %q correctly set as default.\n", name))

			return nil
		},
	}

	return cmd
}
