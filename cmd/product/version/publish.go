package version

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

const (
	_forceFlag = "force"
)

// NewPublishCmd creates a new command to publish a version.
func NewPublishCmd(logger logging.Interface) *cobra.Command {
	nArgs := 3
	cmd := &cobra.Command{
		Use: "publish <product_id> <version_tag> <comment> [opts...]",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Args:  cobra.ExactArgs(nArgs),
		Short: "Publish an started version",
		Example: `
			$ kli product version publish <product_id> <version_tag> <comment> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]
			versionTag := args[1]
			comment := args[2]

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			force, err := cmd.Flags().GetBool(_forceFlag)
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			return version.NewHandler(logger, r, api.NewKaiClient().VersionClient(), api.NewKaiClient().ProductClient()).
				Publish(&version.PublishOpts{
					ServerName: server,
					ProductID:  productID,
					VersionTag: versionTag,
					Comment:    comment,
					Force:      force,
				})
		},
	}

	cmd.Flags().BoolP(_forceFlag, "f", false, "The version of the product.")

	return cmd
}
