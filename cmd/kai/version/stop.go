package version

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/kai"
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
			kaiClient := api.NewKAIClient(&graphql.ClientConfig{
				Debug:                 cfg.Debug,
				DefaultRequestTimeout: cfg.DefaultRequestTimeout,
			}, server, cfg.BuildVersion)
			kaiInteractor := kai.NewInteractor(logger, kaiClient, nil)

			err = kaiInteractor.StopVersion(product, versionName, comment)
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