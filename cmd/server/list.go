package server

import (
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/server"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new command to list servers existing in config file.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List all available servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			return server.NewKaiConfigurator(logger, r).ListServers()
		},
	}

	return cmd
}
