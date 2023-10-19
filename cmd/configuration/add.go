package configuration

import (
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	processconfiguration "github.com/konstellation-io/kli/internal/commands/configuration"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewAddCmd creates a new command to add a new configuration to the given scope.
func NewAddCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <key> <value> [opts...]",
		Args:  cobra.ExactArgs(2), //nolint:gomnd
		Short: "Add a new configuration key-value pair to the given scope",
		Example: `
    $ kli config add <key> <value> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

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
				AddConfiguration(&processconfiguration.AddConfigurationOpts{
					ProductID:  productID,
					WorkflowID: workflowID,
					ProcessID:  processID,
					Scope:      scope,
					Key:        key,
					Value:      value,
				})
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID to register the configuration.")
	cmd.Flags().String(_workflowIDFlag, "", "The workflow ID to register the configuration.")
	cmd.Flags().String(_processIDFlag, "", "The process ID to register the configuration.")
	cmd.Flags().String(_scopeFlag, "process", "The scope to register the configuration.")

	return cmd
}
