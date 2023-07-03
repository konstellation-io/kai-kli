package version

import (
	config2 "github.com/konstellation-io/kli/cmd/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/kre"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewStopCmd manages the stop command.
func NewStopCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop <version-name> -r <product> -m <audit message>",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Stop a version",
		RunE: func(cmd *cobra.Command, args []string) error {
			versionName := args[0]

			product, err := cmd.Flags().GetString("product")
			if err != nil {
				return err
			}

			comment, err := cmd.Flags().GetString("message")
			if err != nil {
				return err
			}

			serverName, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			server := cfg.GetByServerName(serverName)
			kreClient := api.NewKreClient(&graphql.ClientConfig{
				Debug:                 viper.GetBool(config2.DebugKey),
				DefaultRequestTimeout: viper.GetDuration("request_timeout"),
			}, server, viper.GetString(config2.BuildVersionKey))
			kreInteractor := kre.NewInteractor(logger, kreClient, nil)

			err = kreInteractor.StopVersion(product, versionName, comment)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("message", "m", "", "Adds audit message")
	_ = cmd.MarkFlagRequired("message")

	return cmd
}
