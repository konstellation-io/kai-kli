package product

import (
	"github.com/konstellation-io/kli/api"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	_forceBindFlag = "force"
)

// NewBindCmd creates a new command to bind a product.
func NewBindCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind <product_name> [opts...]",
		Args:  cobra.ExactArgs(1),
		Short: "Bind a new product",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Example: `
    $ kli product bind <product_name> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]

			localPath, err := cmd.Flags().GetString(_localPathFlag)
			if err != nil {
				return err
			}

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			force, err := cmd.Flags().GetBool(_forceBindFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			productConfigService := productconfiguration.NewProductConfigService(logger)

			err = product.NewHandler(logger, r, api.NewKaiClient().ProductClient(), productConfigService).
				Bind(&product.BindProductOpts{
					ProductID: productID,
					LocalPath: localPath,
					Server:    server,
					Force:     force,
				})
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