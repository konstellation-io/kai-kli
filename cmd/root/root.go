package root

import (
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/cmd/kre"
	"github.com/konstellation-io/kli/cmd/krt"
	"github.com/konstellation-io/kli/cmd/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/pkg/iostreams"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the base command where all subcommands are added.
func NewRootCmd(
	cfg *config.Config,
	logger logging.Interface,
	io *iostreams.IOStreams,
	version,
	buildDate string,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kli [command] [subcommand] [flags]",
		Short: "Konstellation CLI",
		Long:  `Use Konstellation API from the command line.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			d, err := cmd.Flags().GetBool("debug")
			if d {
				cfg.Debug = true
				logger.SetDebugLevel()
			}

			return err
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetOut(io.Out)
	cmd.SetErr(io.ErrOut)

	// Hide help command. only --help
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	cmd.PersistentFlags().Bool("debug", false, "Set debug mode")

	// Child commands
	cmd.AddCommand(newVersionCmd(version, buildDate))
	cmd.AddCommand(server.NewServerCmd(logger, cfg))
	cmd.AddCommand(kre.NewKRECmd(logger, cfg))
	cmd.AddCommand(krt.NewKRTCmd(logger))

	return cmd
}
