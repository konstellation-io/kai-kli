package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/configuration"
	"github.com/konstellation-io/kli/internal/logging"

	"github.com/konstellation-io/kli/internal/render"
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

			configHandler := configuration.NewKaiConfigHandler(logger)
			kaiConfig, err := configHandler.GetConfiguration()
			if err != nil {
				return err
			}

			err = kaiConfig.SetDefaultServer(name)
			if err != nil {
				return err
			}

			err = configHandler.WriteConfiguration(kaiConfig)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			r.RenderServers(kaiConfig.Servers)

			logger.Success(fmt.Sprintf("Server %q is now default.\n", name))

			return nil
		},
	}

	return cmd
}
