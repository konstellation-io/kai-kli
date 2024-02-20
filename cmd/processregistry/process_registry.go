package processregistry

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_serverFlag      = "server"
	_publicFlag      = "public"
	_productIDFlag   = "product"
	_processNameFlag = "process"
	_versionFlag     = "version"
	_processTypeFlag = "type"
)

// NewProcessRegistryCmd creates a new command to handle 'process-registry' subcommands.
func NewProcessRegistryCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process-registry <command>",
		Short: "Manage process-registry",
		Example: heredoc.Doc(`
			$ kli process-registry ls <product_id> [opts...]
			$ kli process-registry register <process_type> <process_id> [opts...]
		`),
	}

	cmd.PersistentFlags().StringP(_serverFlag, "s", "", "KAI server to use")

	cmd.AddCommand(
		NewListCmd(logger),
		NewRegisterCmd(logger),
		NewDeleteCmd(logger),
	)

	return cmd
}
