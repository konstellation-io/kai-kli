package product

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new command to bind a product.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [opts...]",
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"ls"},
		Short:   "List existing products.",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Example: `
    $ kli product ls [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			productConfigService := productconfiguration.NewProductConfigService(logger)
			kaiClient := api.NewKaiClient()

			err = product.NewHandler(logger, r, kaiClient.ProductClient(), kaiClient.VersionClient(), productConfigService).
				ListProducts(server)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_localPathFlag, "", "The path where the local environment is initialized.")
	cmd.Flags().BoolP(_forceBindFlag, "f", false, "If true, overrides local product configuration if exists.")

	return cmd
}
