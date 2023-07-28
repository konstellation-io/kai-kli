package workflow

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/workflow"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

var (
	ErrDockerfileNotFound = fmt.Errorf("dockerfile not found, the file must be named 'Dockerfile'")
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

			// TODO Get the given product or the default one
			// TODO Get the given server or the default one

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = workflow.NewHandler(logger, r).AddWorkflow(serverName, productID, workflowID, workflowType)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Workflow successfully registered with ID %s.", workflowID))

			return nil
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")

	return cmd
}
