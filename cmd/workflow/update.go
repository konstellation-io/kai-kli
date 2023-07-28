package workflow

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/workflow"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	_workflowTypeFlag = "type"
)

// NewUpdateCmd creates a new command to update a workflow to the given product.
func NewUpdateCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <workflowID> [--workflow-type <workflow_type>] [--product_id <product_id>] [--server <server_name>]",
		Args:    cobra.ExactArgs(1),
		Short:   "Update an existing workflow for the given product on the given server",
		Example: "$ kli update <workflowID> [--workflow-type <workflow_type>] [--product_id <product_id>] [--server <server_name>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]

			workflowType, err := cmd.Flags().GetString(_workflowTypeFlag)
			if err != nil {
				return err
			}

			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				serverName = ""
			}

			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}

			// TODO Get the given product or the default one
			// TODO Get the given server or the default one

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = workflow.NewHandler(logger, r).UpdateWorkflow(serverName, productID, workflowID, workflowType)
			if err != nil {
				return err
			}

			logger.Success("Workflow successfully updated.")

			return nil
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")
	cmd.Flags().String(_workflowTypeFlag, "", "The workflow type to register the process.")

	return cmd
}
