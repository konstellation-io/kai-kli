package process

import (
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/process"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewListCmd creates a new command to list the existing workflows for the given product.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [--server <server_name>] [--product_id <product_id>]  [--workflow-id <workflow_id>]",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		Short:   "List all the processes for the given product workflow on the given server",
		Example: "$ kli process ls [--server <server_name>] [--product_id <product_id>]  [--workflow-id <workflow_id>]",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			err = process.NewHandler(logger, r).ListProcesses(&process.ListProcessOpts{
				ServerName: serverName,
				ProductID:  productID,
				WorkflowID: workflowID,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_serverFlag, "", "The server name from where to list the process.")
	cmd.Flags().String(_productIDFlag, "", "The product ID from where to list the process.")
	cmd.Flags().String(_workflowIDFlag, "", "The workflow ID from where to list the process.")

	return cmd
}
