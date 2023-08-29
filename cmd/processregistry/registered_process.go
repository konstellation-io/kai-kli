package processregistry

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewRegisteredProcessCmd creates a new command to handle 'process-registry' subcommands.
func NewRegisteredProcessCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process-registry <command>",
		Short: "Manage process-registry",
		Example: heredoc.Doc(`
			$ kli process-registry ls <product_id> [--type <process_type>] [opts...]
			$ kli process-registry register <process_type> <process_id> [opts...]
			$ kli process-registry unregister <process_id> [opts...]
		`),
	}

	cmd.AddCommand(
		NewListCmd(logger),
		NewRegisterCmd(logger),
	)

	return cmd
}
