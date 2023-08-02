package workflow

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_productIDFlag = "product"
	_serverFlag    = "server"
)

// NewWorkflowCmd creates a new command to handle 'workflow' subcommands.
func NewWorkflowCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow <command>",
		Short: "Manage workflows",
		Example: heredoc.Doc(`
			$ kli workflow ls [opts..]
			$ kli workflow add <workflowID> <workflowType> [opts...]
			$ kli workflow update <workflowID> [opts...]
			$ kli workflow remove <workflowID> [opts...]
		`),
	}

	cmd.AddCommand(
		NewAddCmd(logger),
		NewUpdateCmd(logger),
		NewRemoveCmd(logger),
		NewListCmd(logger),
	)

	return cmd
}
