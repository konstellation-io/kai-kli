package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

// NewAddCmd creates a new command to add a new server to config file.
func NewAddCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <server-name> <server-url> <token>",
		Aliases: []string{"set"},
		Args:    cobra.ExactArgs(3), //nolint:gomnd
		Short:   "Add a new server config file",
		Example: `
    $ kli server add my-server http://api.local.kai 12345abc
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			newServer := config.ServerConfig{
				Name:     args[0],
				URL:      args[1],
				APIToken: args[2],
			}

			err := cfg.AddServer(newServer)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			r.RenderServerList(cfg.ServerList, cfg.DefaultServer)

			logger.Success(fmt.Sprintf("Server %q added.", newServer.Name))
			return nil
		},
	}

	return cmd
}
