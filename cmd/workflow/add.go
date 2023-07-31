package workflow

import (
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/workflow"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	_productIDFlag = "product-id"
	_serverFlag    = "server"
)

// NewAddCmd creates a new command to add a workflow to the given product.
func NewAddCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <workflow_id> <workflow_type> [--product-id <product_id>] [--server <server_name>]",
		Args:    cobra.ExactArgs(2), //nolint:gomnd
		Short:   "Registers a new workflow for the given product on the given server",
		Example: "$ kli workflow add <workflow_id> <workflow_type> [--product_id <product_id>] [--server <server_name>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			workflowType := args[1]

			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}

			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				serverName = ""
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			return workflow.NewHandler(logger, r).AddWorkflow(&workflow.AddWorkflowOpts{
				ServerName:   serverName,
				ProductID:    productID,
				WorkflowID:   workflowID,
				WorkflowType: krt.WorkflowType(workflowType),
			})
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")

	return cmd
}
