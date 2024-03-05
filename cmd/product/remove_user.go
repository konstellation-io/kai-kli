package product

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/render"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewRemoveUserCmd creates a new command to remove an user from a product.
func NewRemoveUserCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-user <product_name> <user_email> [opts...]",
		Args:  cobra.ExactArgs(2), //nolint:gomnd
		Short: "Remove user from product",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Example: `
    $ kli product remove-user <product_name> <user_email> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]
			userEmail := args[1]

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			productConfigService := productconfiguration.NewProductConfigService(logger)
			kaiClient := api.NewKaiClient()

			err = product.NewHandler(logger, r, kaiClient.ProductClient(), kaiClient.VersionClient(), productConfigService).
				RemoveUser(server, productID, userEmail)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
