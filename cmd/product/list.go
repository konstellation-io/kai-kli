package product

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/spf13/cobra"
)

const (
	_productFlag = "product"
)

// NewListCmd creates a new command to bind a product.
func NewListCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [--product <product_name>] [opts...]",
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"ls"},
		Short:   "List existing products.",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Example: `
    $ kli product ls [opts...]
		`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var opts product.ListProductsOpts

			productName, err := cmd.Flags().GetString(_productFlag)
			if err == nil {
				opts.ProductName = productName
			}

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			productConfigService := productconfiguration.NewProductConfigService(logger)
			kaiClient := api.NewKaiClient()

			err = product.NewHandler(
				logger, r, kaiClient.ProductClient(), kaiClient.VersionClient(), productConfigService,
			).ListProducts(server, &opts)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_productFlag, "", "Filter results by product name.")

	return cmd
}
