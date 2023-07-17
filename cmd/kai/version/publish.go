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

// NewPublishCmd manages the publish command.
func NewPublishCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish <version-name> -r <product> -m <audit message>",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Publish a version",
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
			kaiClient := api.NewKaiClient(&graphql.ClientConfig{
				Debug:                 viper.GetBool(config2.DebugKey),
				DefaultRequestTimeout: viper.GetDuration("request_timeout"),
			}, server, viper.GetString(config2.BuildVersionKey))
			kaiInteractor := kai.NewInteractor(logger, kaiClient, nil)

			err = kaiInteractor.PublishVersion(product, versionName, comment)
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
