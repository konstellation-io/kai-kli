package processregistry

import (
	"fmt"
	"strings"

	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/processregistry"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

var (
	ErrDockerfileNotFound = fmt.Errorf("dockerfile not found, the file must be named 'Dockerfile'")
)

const (
	_callbackPath   = "sso-callback"
	_sourcesFlag    = "src"
	_dockerfileFlag = "dockerfile"
	_productIDFlag  = "product"
	_serverFlag     = "server"
	_versionFlag    = "version"
)

// NewRegisterCmd creates a new command to register a new process in the given server.
func NewRegisterCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "register <process_type> <process_id> [opts...]",
		Args: cobra.ExactArgs(2), //nolint:gomnd
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Short:   "Register a new process for the given product on the given server",
		Example: "$ kli process-registry register <process_type> <process_id> [opts...]",
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
			err = processregistry.NewHandler(
				logger,
				r,
				api.NewKaiClient().ProcessRegistry(),
			).
				RegisterProcess(&processregistry.RegisterProcessOpts{
					ServerName:  serverName,
					ProductID:   productID,
					ProcessType: krt.ProcessType(processType),
					ProcessID:   processID,
					SourcesPath: sourcesPath,
					Dockerfile:  dockerFile,
					Version:     version,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_sourcesFlag, "", "Path to the source code of the process.")
	cmd.Flags().String(_dockerfileFlag, "", "Path to the Dockerfile of the process.")
	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")
	cmd.Flags().String(_versionFlag, "s", "The version with which register the process.")
	cmd.Flags().String(_serverFlag, "v", "The server where the process will be registered.")

	cmd.MarkFlagRequired(_sourcesFlag)    //nolint:errcheck
	cmd.MarkFlagRequired(_dockerfileFlag) //nolint:errcheck
	cmd.MarkFlagRequired(_versionFlag)    //nolint:errcheck
	cmd.MarkFlagRequired(_productIDFlag)  //nolint:errcheck

	return cmd
}
