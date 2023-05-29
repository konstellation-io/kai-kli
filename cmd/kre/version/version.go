package version

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kre/config"
	configcmd "github.com/konstellation-io/kli/cmd/kre/version/config"
	"github.com/konstellation-io/kli/internal/logging"
)

const productFlag = "product"

// NewVersionCmd creates a new command to handle 'version' subcommands.
func NewVersionCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Manage KRE Version",
		Example: heredoc.Doc(`
			$ kli kre version ls -r demo
			$ kli kre version create file.krt -r demo
			$ kli kre version start demo-v1 -r demo -m 'Testing version'
		`),
	}

	cmd.PersistentFlags().StringP("server", "s", cfg.DefaultServer, "KRE server to use")
	cmd.PersistentFlags().StringP(productFlag, "r", "", "product of the version")
	_ = cmd.MarkFlagRequired(productFlag)

	cmd.AddCommand(
		NewListCmd(logger, cfg),
		NewStartCmd(logger, cfg),
		NewStopCmd(logger, cfg),
		NewPublishCmd(logger, cfg),
		NewUnpublishCmd(logger, cfg),
		NewCreateCmd(logger, cfg),
		configcmd.NewConfigCmd(logger, cfg),
	)

	return cmd
}
