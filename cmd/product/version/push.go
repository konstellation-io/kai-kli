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
)

// NewPushCmd creates a new command to push a new product's version.
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

			versionNum, err := cmd.Flags().GetString(_versionFlag)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(_descriptionFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = version.NewHandler(logger, r, api.NewKaiClient().VersionClient(), api.NewKaiClient().ProductClient()).
				PushVersion(&version.PushVersionOpts{
					Server:      server,
					KrtFilePath: krtFilePath,
					Version:     versionNum,
					Description: description,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_versionFlag, "", "The version of the product.")
	cmd.Flags().String(_descriptionFlag, "", "The description of the product.")

	return cmd
}
