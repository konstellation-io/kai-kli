package configuration

import (
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	processconfiguration "github.com/konstellation-io/kli/internal/commands/configuration"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewListCmd creates a new command to list existing configuration from the given scope.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [opts...]",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		Short:   "List the existing configuration from the given scope",
		Example: `
    $ kli config ls [opts...]
		`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}

			workflowID, err := cmd.Flags().GetString(_workflowIDFlag)
			if err != nil {
				return err
			}

			processID, err := cmd.Flags().GetString(_processIDFlag)
			if err != nil {
				return err
			}

			scope, err := cmd.Flags().GetString(_scopeFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			return processconfiguration.NewHandler(logger, r).
				ListConfiguration(&processconfiguration.ListConfigurationOpts{
					ProductID:  productID,
					WorkflowID: workflowID,
					ProcessID:  processID,
					Scope:      scope,
				})
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID to register the configuration.")
	cmd.Flags().String(_workflowIDFlag, "", "The workflow ID to register the configuration.")
	cmd.Flags().String(_processIDFlag, "", "The process ID to register the configuration.")
	cmd.Flags().String(_scopeFlag, "process", "The scope to register the configuration.")

	return cmd
}
