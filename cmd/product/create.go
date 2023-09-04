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
	_versionFlag     = "version"
	_descriptionFlag = "description"
	_initLocalFlag   = "init-local"
	_localPathFlag   = "local-path"
)

// NewCreateCmd creates a new command to create a new product.
func NewCreateCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <product_name> [opts...]",
		Aliases: []string{"add"},
		Args:    cobra.ExactArgs(1),
		Short:   "Add a new product",
		Example: `
    $ kli product create <product_name> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productName := args[0]

			version, err := cmd.Flags().GetString(_versionFlag)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(_descriptionFlag)
			if err != nil {
				return err
			}

			initLocal, err := cmd.Flags().GetBool(_initLocalFlag)
			if err != nil {
				return err
			}

			localPath, err := cmd.Flags().GetString(_localPathFlag)
			if err != nil {
				return err
			}

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			productConfigService := productconfiguration.NewProductConfigService(logger)

			err = product.NewHandler(logger, r, api.NewKaiClient().ProductClient(), productConfigService).CreateProduct(&product.CreateProductOpts{
				ProductName: productName,
				Version:     version,
				Description: description,
				InitLocal:   initLocal,
				LocalPath:   localPath,
				Server:      server,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_versionFlag, "v0.0.1", "The version of the product.")
	cmd.Flags().String(_descriptionFlag, "", "The description of the product.")
	cmd.Flags().Bool(_initLocalFlag, false, "If true, a local product environment is initialized.")
	cmd.Flags().String(_localPathFlag, "", "The path where the local environment is initialized.")

	return cmd
}
