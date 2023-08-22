package processregistry

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api"
	processregistry "github.com/konstellation-io/kli/internal/commands/process-registry"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/auth"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

var (
	ErrDockerfileNotFound = fmt.Errorf("dockerfile not found, the file must be named 'Dockerfile'")
)

const (
	_sourcesFlag    = "src"
	_dockerfileFlag = "dockerfile"
	_productIDFlag  = "product-id"
	_serverFlag     = "server"
	_versionFlag    = "version"
)

// NewRegisterCmd creates a new command to register a new process in the given server.
func NewRegisterCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use: "register <process_type> <process_id>" +
			" --src <sources> --dockerfile <dockerfile> --product_id <product_id> [--server <server_name>]",
		Args: cobra.ExactArgs(2), //nolint:gomnd
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Short: "Register a new process for the given product on the given server",
		Example: "$ kli process-registry register <process_type> <process_id> --src <sources>" +
			" --dockerfile <dockerfile> --product-id <product_id> [--server <server_name>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			processType := args[0]
			processID := args[1]

			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}
			sourcesPath, err := cmd.Flags().GetString(_sourcesFlag)
			if err != nil {
				return err
			}
			dockerFile, err := cmd.Flags().GetString(_dockerfileFlag)
			if err != nil {
				return err
			}
			version, err := cmd.Flags().GetString(_versionFlag)
			if err != nil {
				return err
			}

			if !strings.HasSuffix(dockerFile, "Dockerfile") {
				return ErrDockerfileNotFound
			}

			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				serverName = ""
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = processregistry.NewProcessRegistryHandler(
				logger,
				r,
				api.NewKaiClient().ProcessRegistry(),
				auth.NewAuthentication(logger),
				configuration.NewKaiConfigService(logger),
			).
				RegisterProcess(serverName, productID, processType,
					processID, sourcesPath, dockerFile, version)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_sourcesFlag, "", "Path to the source code of the process.")
	cmd.Flags().String(_dockerfileFlag, "", "Path to the Dockerfile of the process.")
	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")
	cmd.Flags().StringP(_serverFlag, "s", "", "The product ID to register the process.")
	cmd.Flags().StringP(_versionFlag, "v", "", "The product ID to register the process.")

	cmd.MarkFlagRequired(_sourcesFlag)    //nolint:errcheck
	cmd.MarkFlagRequired(_dockerfileFlag) //nolint:errcheck
	cmd.MarkFlagRequired(_productIDFlag)  //nolint:errcheck

	return cmd
}
