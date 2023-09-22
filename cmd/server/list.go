package server

import (
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewListCmd creates a new command to list servers existing in config file.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List all available servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			return server.NewHandler(logger, r).ListServers()
		},
	}

	return cmd
}
