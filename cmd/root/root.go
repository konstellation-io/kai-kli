package root

import (
	"github.com/konstellation-io/kli/authserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/konstellation-io/kli/cmd/config"
	configuration2 "github.com/konstellation-io/kli/cmd/configuration"
	"github.com/konstellation-io/kli/cmd/process"
	"github.com/konstellation-io/kli/cmd/processregistry"
	"github.com/konstellation-io/kli/cmd/product"
	"github.com/konstellation-io/kli/cmd/server"
	"github.com/konstellation-io/kli/cmd/workflow"
	"github.com/konstellation-io/kli/internal/commands/setup"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/auth"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/pkg/iostreams"
)

const (
	_serverFlag                     = "server"
	_debugFlag                      = "debug"
	_authenticatedAnnotation        = "authenticated"
	_productConfigurationAnnotation = "needs-product-config"
)

// NewRootCmd creates the base command where all subcommands are added.
func NewRootCmd(
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
			err := setDebugLogLevel(cmd, logger)
			if err != nil {
				return err
			}

			err = setup.NewHandler(logger).CreateConfiguration()
			if err != nil {
				return err
			}

			if isMethodAuthenticated(cmd) {
				err = authenticateServer(logger, cmd)
				if err != nil {
					return err
				}
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

	cmd.PersistentFlags().Bool(_debugFlag, false, "Set debug mode")

	// Child commands
	cmd.AddCommand(newVersionCmd(version, buildDate))
	cmd.AddCommand(server.NewServerCmd(logger))
	cmd.AddCommand(product.NewProductCmd(logger))
	cmd.AddCommand(workflow.NewWorkflowCmd(logger))
	cmd.AddCommand(process.NewProcessCmd(logger))
	cmd.AddCommand(configuration2.NewConfigurationCmd(logger))
	cmd.AddCommand(processregistry.NewProcessRegistryCmd(logger))

	return cmd
}

func authenticateServer(logger logging.Interface, cmd *cobra.Command) error {
	configService := configuration.NewKaiConfigService(logger)
	kaiAuth := auth.NewAuthentication(logger, authserver.NewDefaultAuthServer(logger))

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	serverName, _ := cmd.Flags().GetString(_serverFlag)

	srv, err := conf.GetServerOrDefault(serverName)
	if err != nil {
		return err
	}

	if _, err := kaiAuth.GetToken(srv.Name); err != nil {
		return err
	}

	return nil
}

func isMethodAuthenticated(cmd *cobra.Command) bool {
	val, ok := cmd.Annotations[_authenticatedAnnotation]
	return ok && val == "true"
}

func setDebugLogLevel(cmd *cobra.Command, logger logging.Interface) error {
	d, err := cmd.Flags().GetBool(_debugFlag)

	if d {
		viper.Set(config.DebugKey, true)

		logger.SetDebugLevel()
	}

	return err
}
