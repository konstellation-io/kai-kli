package workflow

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/workflow"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

// NewListCmd creates a new command to list the existing workflows for the given product.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [--product-id <product_id>] [--server <server_name>]",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		Short:   "List all the workflows for the given product on the given server",
		Example: "$ kli ls [--product_id <product_id>] [--server <server_name>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			return workflow.NewHandler(logger, r).ListWorkflows(&workflow.ListWorkflowOpts{
				ProductID: productID,
			})
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")

	return cmd
}
