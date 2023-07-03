package product

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/cmd/args"
	config2 "github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/kre"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListCmd creates a new command to list Products.
func NewListCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Args:    args.CheckServerFlag,
		Short:   "List all available Products",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, _ := cmd.Flags().GetString("server")
			server := cfg.GetByServerName(serverName)
			kreClient := api.NewKreClient(&graphql.ClientConfig{
				Debug:                 viper.GetBool(config2.DebugKey),
				DefaultRequestTimeout: viper.GetDuration("request_timeout"),
			}, server, viper.GetString(config2.BuildVersionKey))

			kreInteractor := kre.NewInteractorWithDefaultRenderer(logger, kreClient, cmd.OutOrStdout())
			return kreInteractor.ListProducts()
		},
	}

	return cmd
}
