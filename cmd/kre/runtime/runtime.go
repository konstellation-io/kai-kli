package runtime

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewRuntimeCmd creates a new command to handle 'runtime' subcommands.
func NewRuntimeCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runtime",
		Short: "Manage KRE Version",
		Example: heredoc.Doc(`
			$ kli kre runtime ls
			$ kli kre runtime create 'runtime-name' -d 'runtime description'
		`),
	}

	cmd.PersistentFlags().StringP("server", "s", cfg.DefaultServer, "KRE server to use")

	cmd.AddCommand(
		NewListCmd(logger, cfg),
		NewCreateCmd(logger, cfg),
	)

	return cmd
}
