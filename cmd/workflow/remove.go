package workflow

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/workflow"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

// NewRemoveCmd creates a new command to delete a workflow to the given product.
func NewRemoveCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <workflowID> [--product-id <product_id>] [--server <server_name>]",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"rm", "del", "delete"},
		Short:   "Remove an existing workflow for the given product on the given server",
		Example: "$ kli workflow remove <workflowID> [--product-id <product_id>] [--server <server_name>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]

			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			return workflow.NewHandler(logger, r).RemoveWorkflow(&workflow.RemoveWorkflowOpts{
				ProductID:  productID,
				WorkflowID: workflowID,
			})
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")

	return cmd
}
