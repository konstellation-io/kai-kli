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

// NewCreateCmd upload a KRT file to product and make a new version.
func NewCreateCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [krt-file] -r <product-id>",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Upload a KRT and create a new version",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, _ := cmd.Flags().GetString("server")
			server := cfg.GetByServerName(serverName)
			kreClient := api.NewKreClient(&graphql.ClientConfig{
				Debug:                 cfg.Debug,
				DefaultRequestTimeout: cfg.DefaultRequestTimeout,
			}, server, cfg.BuildVersion)
			kreInteractor := kre.NewInteractor(logger, kreClient, nil)

			krt := args[0]

			product, _ := cmd.Flags().GetString(productFlag)

			return kreInteractor.CreateVersion(product, krt)
		},
	}

	return cmd
}
