package version

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates a new command to handle 'version' subcommands.
func NewVersionCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version <command>",
		Aliases: []string{"version", "versions"},
		Short:   "Manage versions",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Example: heredoc.Doc(`
			$ kli version ls [opts...]
			$ kli version create <product_name> [opts...]
			$ kli version bind <product_name> [opts...]
		`),
	}

	cmd.PersistentFlags().StringP("server", "s", "", "KAI server to use")

	cmd.AddCommand(
		NewPushCmd(logger),
		NewStartCmd(logger),
		NewStopCmd(logger),
	)

	return cmd
}
