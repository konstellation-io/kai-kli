package kre

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/cmd/kre/runtime"
	"github.com/konstellation-io/kli/cmd/kre/server"
	"github.com/konstellation-io/kli/cmd/kre/version"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewKRECmd creates a new command to handle 'kre' keyword.
func NewKRECmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kre",
		Short: "Manage KRE",
		Example: heredoc.Doc(`
			$ kli kre runtime ls
		`),
	}

	cmd.AddCommand(
		runtime.NewRuntimeCmd(logger, cfg),
		version.NewVersionCmd(logger, cfg),
		server.NewServerCmd(logger, cfg),
	)

	return cmd
}
