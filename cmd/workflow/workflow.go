package workflow

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewWorkflowCmd creates a new command to handle 'workflow' subcommands.
func NewWorkflowCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow <command>",
		Short: "Manage workflows",
		Example: heredoc.Doc(`
			$ kli workflow ls [--product-id <product_id>] [--server <server_name>]
			$ kli workflow add <workflowID> <workflowType> [--product_id <product_id>] [--server <server_name>]
			$ kli workflow update <workflowID> [--workflow-type <workflow_type>] [--product_id <product_id>] [--server <server_name>]
			$ kli workflow remove <workflowID> [--product_id <product_id>] [--server <server_name>]
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
