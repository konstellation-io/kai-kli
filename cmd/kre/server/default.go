package server

import (
	"fmt"

	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/render"
)

// NewDefaultCmd creates a new command to handle 'default' keyword.
func NewDefaultCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default <server-name>",
		Short: "Set a default server",
		Long:  "Mark an existing server as default",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			err := cfg.SetDefaultServer(name)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			r.RenderServerList(cfg.ServerList, cfg.DefaultServer)

			logger.Success(fmt.Sprintf("Server '%s' is now default.\n", name))

			return nil
		},
	}

	return cmd
}
