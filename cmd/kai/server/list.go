package server

import (
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/render"
)

// NewListCmd creates a new command to list servers existing in config file.
func NewListCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List all available servers",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg.ServerList) == 0 {
				logger.Info("No servers found.")
				return
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			r.RenderServerList(cfg.ServerList, cfg.DefaultServer)
		},
	}

	return cmd
}
