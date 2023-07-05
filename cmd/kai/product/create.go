package product

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/kai"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates a new command to list Products.
func NewCreateCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <product-name> -d <product description>",
		Aliases: []string{"create"},
		Args:    args.CheckServerFlag,
		Short:   "Create a product",
		Example: heredoc.Doc(`
			$ kli kai product create 'product-name' -d 'product description' -s kai-local
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, _ := cmd.Flags().GetString("server")
			server := cfg.GetByServerName(serverName)
			kaiClient := api.NewKAIClient(&graphql.ClientConfig{
				Debug:                 cfg.Debug,
				DefaultRequestTimeout: cfg.DefaultRequestTimeout,
			}, server, cfg.BuildVersion)

			name := args[0]
			description, _ := cmd.Flags().GetString("description")

			kaiInteractor := kai.NewInteractorWithDefaultRenderer(logger, kaiClient, cmd.OutOrStdout())
			err := kaiInteractor.CreateProduct(name, description)
			if err != nil {
				logger.Debug(err.Error())
				return err
			}

			logger.Info("Product successfully created")
			return nil
		},
	}
	cmd.Flags().StringP("description", "d", "", "Adds product description")

	return cmd
}
