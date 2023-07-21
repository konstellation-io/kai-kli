package processregistry

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewProcessRegistryCmd creates a new command to handle 'process-registry' subcommands.
func NewProcessRegistryCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process-registry <command>",
		Short: "Manage process-registry",
		Example: heredoc.Doc(`
			$ kli process-registry ls [--product-id <product_id>] [--process-type <process_type>]
			$ kli process-registry register <process_type> <process_id> [--src <sources>] [--dockerfile <dockerfile>] [--product_id <product_id>]
			$ kli process-registry unregister <process_id> [--product_id <product_id>]
		`),
	}

	cmd.AddCommand(
		NewRegisterCmd(logger),
	)

	return cmd
}
