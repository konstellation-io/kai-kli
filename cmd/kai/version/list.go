package version

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	config2 "github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/kai"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewListCmd creates a new command to list Versions.
func NewListCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls -r <product>",
		Aliases: []string{"list"},
		Args:    args.CheckServerFlag,
		Short:   "List all available Versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, _ := cmd.Flags().GetString("server")
			server := cfg.GetByServerName(serverName)
			kaiClient := api.NewKaiClient(&graphql.ClientConfig{
				Debug:                 viper.GetBool(config2.DebugKey),
				DefaultRequestTimeout: viper.GetDuration("request_timeout"),
			}, server, viper.GetString(config2.BuildVersionKey))
			kaiInteractor := kai.NewInteractorWithDefaultRenderer(logger, kaiClient, cmd.OutOrStdout())

			product, _ := cmd.Flags().GetString(productFlag)

			return kaiInteractor.ListVersions(product)
		},
	}

	return cmd
}
