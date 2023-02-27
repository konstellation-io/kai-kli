package version

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/kre"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewPublishCmd manages the publish command.
func NewPublishCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish <version-name> -r <runtime> -m <audit message>",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Publish a version",
		RunE: func(cmd *cobra.Command, args []string) error {
			versionName := args[0]

			runtime, err := cmd.Flags().GetString("runtime")
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
				Debug:                 cfg.Debug,
				DefaultRequestTimeout: cfg.DefaultRequestTimeout,
			}, server, cfg.BuildVersion)
			kreInteractor := kre.NewInteractor(logger, kreClient, nil)

			err = kreInteractor.PublishVersion(runtime, versionName, comment)
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
