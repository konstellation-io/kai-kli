package version

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

const (
	_versionFlag     = "version"
	_descriptionFlag = "description"
	_initLocalFlag   = "init-local"
	_localPathFlag   = "local-path"
)

// NewPushCmd creates a new command to create a new product.
func NewPushCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use: "push <krt_file> [opts...]",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Push a new version",
		Example: `
    $ kli product version push <krt_file> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			krtFilePath := args[0]

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = version.NewHandler(logger, r, api.NewKaiClient().VersionClient()).
				PushVersion(&version.PushVersionOpts{
					Server:      server,
					KrtFilePath: krtFilePath,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_versionFlag, "v0.0.1", "The version of the product.")
	cmd.Flags().String(_descriptionFlag, "", "The description of the product.")
	cmd.Flags().Bool(_initLocalFlag, false, "If true, a local product environment is initialized.")
	cmd.Flags().String(_localPathFlag, "", "The path where the local environment is initialized.")

	return cmd
}
