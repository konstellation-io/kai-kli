package process

import (
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/process"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewRemoveCmd creates a new command to delete a workflow to the given product.
func NewRemoveCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <process_id> [--server <server_name>] [--product-id <product_id>] [--workflow-id <workflow_id>]",
		Aliases: []string{"rm", "del", "delete"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove an existing process for the given product workflow on the given server",
		Example: "$ kli process remove <process_id> [--server <server_name>] [--product-id <product_id>] [--workflow-id <workflow_id>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			processID := args[0]

			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				serverName = ""
			}

			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}

			workflowID, err := cmd.Flags().GetString(_workflowIDFlag)
			if err != nil {
				return err
			}

			// TODO Get the given product or the default one
			// TODO Get the given server or the default one

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = process.NewHandler(logger, r).RemoveProcess(&process.RemoveProcessOpts{
				ServerName: serverName,
				ProductID:  productID,
				WorkflowID: workflowID,
				ProcessID:  processID,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_serverFlag, "", "The product ID to remove the process from.")
	cmd.Flags().String(_productIDFlag, "", "The product ID to remove the process from.")
	cmd.Flags().String(_workflowIDFlag, "", "The workflow ID to remove the process from.")

	return cmd
}
