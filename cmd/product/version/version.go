package version

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates a new command to handle 'version' subcommands.
func NewVersionCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version <command>",
		Aliases: []string{"version", "versions"},
		Short:   "Manage versions",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Example: heredoc.Doc(`
			$ kli product versions ls <product_name> [opts...]
			$ kli product version get <product_name> <version_tag> [opts...]
			$ kli product version publish <product_id> <version_tag> <comment> [opts...]
			$ kli product version unpublish <product_id> <version_tag> <comment> [opts...]
			$ kli product version start <product_id> <version_tag> <comment> [opts...]
			$ kli product version stop <product_id> <version_tag> <comment> [opts...]
			$ kli product version logs <product_id> <version_tag> --from <from_timestamp> --to <to_timestamp> [opts...]
			$ kli product version push <krt_file> [opts...]
		`),
	}

	cmd.PersistentFlags().StringP("server", "s", "", "KAI server to use")

	cmd.AddCommand(
		NewPushCmd(logger),
		NewStartCmd(logger),
		NewStopCmd(logger),
		NewPublishCmd(logger),
		NewGetCmd(logger),
		NewListCmd(logger),
		NewUnpublishCmd(logger),
		NewLogsCmd(logger),
	)

	return cmd
}
