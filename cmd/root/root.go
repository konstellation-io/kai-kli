package root

import (
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/kai"
	"github.com/konstellation-io/kli/cmd/krt"
	"github.com/konstellation-io/kli/cmd/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/setup"
	"github.com/konstellation-io/kli/pkg/iostreams"
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
			err := setDebugLogLevel(cmd, cfg, logger)
			if err != nil {
				return err
			}

			err = setup.NewKaiSetup(logger).CreateConfiguration()
			if err != nil {
				return err
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
	cmd.AddCommand(server.NewServerCmd(logger))
	cmd.AddCommand(kai.NewKAICmd(logger, cfg))
	cmd.AddCommand(krt.NewKRTCmd(logger))

	return cmd
}

func setDebugLogLevel(cmd *cobra.Command, cfg *config.Config, logger logging.Interface) error {
	d, err := cmd.Flags().GetBool("debug")

	if d {
		cfg.Debug = true

		logger.SetDebugLevel()
	}

	return err
}
