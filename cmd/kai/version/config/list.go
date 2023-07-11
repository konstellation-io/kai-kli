package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/args"
	config2 "github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/kai"
	"github.com/konstellation-io/kli/internal/logging"
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
			kaiClient := api.NewKaiClient(&graphql.ClientConfig{
				Debug:                 viper.GetBool(config2.DebugKey),
				DefaultRequestTimeout: viper.GetDuration("request_timeout"),
			}, server, viper.GetString(config2.BuildVersionKey))
			kaiInteractor := kai.NewInteractorWithDefaultRenderer(logger, kaiClient, cmd.OutOrStdout())

			return kaiInteractor.ListVersionConfig(product, versionName, showValues)
		},
	}
	cmd.Flags().Bool("show-values", false, "Show configuration variables")

	return cmd
}
