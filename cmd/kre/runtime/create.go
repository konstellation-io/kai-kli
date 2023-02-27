package runtime

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/cmd/args"
	"github.com/konstellation-io/kli/internal/kre"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates a new command to list Runtimes.
func NewCreateCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <runtime-name> -d <runtime description>",
		Aliases: []string{"create"},
		Args:    args.CheckServerFlag,
		Short:   "Create a runtime",
		Example: heredoc.Doc(`
			$ kli kre runtime create 'runtime-name' -d 'runtime description' -s kre-local
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, _ := cmd.Flags().GetString("server")
			server := cfg.GetByServerName(serverName)
			kreClient := api.NewKreClient(&graphql.ClientConfig{
				Debug:                 cfg.Debug,
				DefaultRequestTimeout: cfg.DefaultRequestTimeout,
			}, server, cfg.BuildVersion)

			name := args[0]
			description, _ := cmd.Flags().GetString("description")

			kreInteractor := kre.NewInteractorWithDefaultRenderer(logger, kreClient, cmd.OutOrStdout())
			err := kreInteractor.CreateRuntime(name, description)
			if err != nil {
				logger.Debug(err.Error())
				return err
			}

			logger.Info("Runtime successfully created")
			return nil
		},
	}
	cmd.Flags().StringP("description", "d", "", "Adds runtime description")

	return cmd
}
