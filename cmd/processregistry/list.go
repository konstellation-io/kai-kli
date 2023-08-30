package processregistry

import (
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/processregistry"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	_processTypeFlag = "type"
)

// NewListCmd creates a new command to list registered process for a given product.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <product_id> [--type <process_type>] [opts...]",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Short:   "List all registered processes for a given product with optional filters",
		Example: "$ kli process-registry ls <product_id> [opts...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]

			processType, err := cmd.Flags().GetString(_processTypeFlag)
			if err != nil {
				processType = ""
			}

			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				serverName = ""
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = processregistry.NewHandler(
				logger,
				r,
				api.NewKaiClient().ProcessRegistry(),
			).
				ListProcesses(&processregistry.ListProcessesOpts{
					ServerName:  serverName,
					ProductID:   productID,
					ProcessType: krt.ProcessType(processType),
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_processTypeFlag, "", "Optional filter results by process type.")

	return cmd
}
