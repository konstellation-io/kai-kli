package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

// newVersionCmd creates a new command to handle 'version' keyword.
func newVersionCmd(logger logging.Interface, version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of the CLI",
		Example: heredoc.Doc(`
			$ kli version
		`),
		Run: func(cmd *cobra.Command, args []string) {
			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			r.RenderKliVersion(version, buildDate)
		},
	}

	return cmd
}
