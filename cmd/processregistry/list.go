package processregistry

import (
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/processregistry"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewListCmd creates a new command to list registered process for a given product.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <product_id> [--process <process_name>] [--version <process_version>] [--type <process_type>] [opts...]",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Short:   "List all registered processes for a given product with optional filters",
		Example: "$ kli process-registry ls <product_id> [opts...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]

			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				serverName = ""
			}

			listFilters := getListFilters(cmd)

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = processregistry.NewHandler(
				logger,
				r,
				api.NewKaiClient().ProcessRegistry(),
			).
				ListProcesses(&processregistry.ListProcessesOpts{
					ServerName: serverName,
					ProductID:  productID,
					Filters:    listFilters,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_processNameFlag, "", "Optional filter results by process name.")
	cmd.Flags().String(_versionFlag, "", "Optional filter results by process version.")
	cmd.Flags().String(_processTypeFlag, "", "Optional filter results by process type.")

	return cmd
}

func getListFilters(cmd *cobra.Command) *processregistry.ListFilters {
	var listFilters processregistry.ListFilters

	processName, err := cmd.Flags().GetString(_processNameFlag)
	if err == nil {
		listFilters.ProcessName = processName
	}

	version, err := cmd.Flags().GetString(_versionFlag)
	if err == nil {
		listFilters.Version = version
	}

	processType, err := cmd.Flags().GetString(_processTypeFlag)
	if err == nil {
		listFilters.ProcessType = krt.ProcessType(processType)
	}

	return &listFilters
}
