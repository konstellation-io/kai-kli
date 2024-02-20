package processregistry

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/processregistry"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewDeleteCmd creates a new command to delete a registered process in the given server.
func NewDeleteCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete <process_id> <version> [--product <product>] [--public]  [opts...]",
		Args: cobra.ExactArgs(2), //nolint:gomnd
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Short: "Delete a registered process for the given product on the given server",
		Example: heredoc.Doc(`
		$ kli process-registry delete <process_id> <version> --product <product> [opts...]
		$ kli process-registry delete <process_id> <version> --public [opts...]
	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			processID := args[0]
			version := args[1]

			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				productID = ""
			}

			isPublic, err := cmd.Flags().GetBool(_publicFlag)
			if err != nil {
				isPublic = false
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
				DeleteProcess(&processregistry.DeleteProcessOpts{
					ServerName: serverName,
					ProductID:  productID,
					ProcessID:  processID,
					Version:    version,
					IsPublic:   isPublic,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_productIDFlag, "", "The product ID process is registered to.")
	cmd.Flags().Bool(_publicFlag, false, "Process is stored in public repo.")

	return cmd
}
