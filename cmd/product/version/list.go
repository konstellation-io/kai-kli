package version

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

const (
	_statusFlag = "status"
)

// NewListCmd creates a new command to list versions of a product.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <product_id> [opts...]",
		Aliases: []string{"ls"},
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Args:  cobra.ExactArgs(1),
		Short: "List all versions of a product",
		Example: `
    	$ kli product version ls <product_id> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			statusFilter, err := cmd.Flags().GetString(_statusFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = version.NewHandler(logger, r, api.NewKaiClient().VersionClient(), api.NewKaiClient().ProductClient()).
				ListVersions(&version.ListVersionsOpts{
					ServerName:   server,
					ProductID:    productID,
					StatusFilter: statusFilter,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_statusFlag, "", "Filter results by status (started, stopped, published, critical).")

	return cmd
}
