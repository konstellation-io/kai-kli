package config

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/kai"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewListConfigCmd list version's configuration variables.
func NewListConfigCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Get config variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			versionName := args[0]

			product, _ := cmd.Flags().GetString("product")

			showValues, err := cmd.Flags().GetBool("show-values")
			if err != nil {
				return err
			}

			serverName, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			server := cfg.GetByServerName(serverName)
			kaiClient := api.NewKAIClient(&graphql.ClientConfig{
				Debug:                 cfg.Debug,
				DefaultRequestTimeout: cfg.DefaultRequestTimeout,
			}, server, cfg.BuildVersion)
			kaiInteractor := kai.NewInteractorWithDefaultRenderer(logger, kaiClient, cmd.OutOrStdout())

			return kaiInteractor.ListVersionConfig(product, versionName, showValues)
		},
	}
	cmd.Flags().Bool("show-values", false, "Show configuration variables")

	return cmd
}