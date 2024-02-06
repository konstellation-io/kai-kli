package version

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

// NewStopCmd creates a new command to stop an existing version.
func NewStopCmd(logger logging.Interface) *cobra.Command {
	nArgs := 3
	cmd := &cobra.Command{
		Use: "stop <product_id> <version_tag> <comment> [opts...]",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Args:  cobra.ExactArgs(nArgs),
		Short: "Stop an existing version",
		Example: `
    	$ kli product version stop <product_id> <version_tag> <comment> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]
			versionTag := args[1]
			comment := args[2]

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = version.NewHandler(logger, r, api.NewKaiClient().VersionClient(), api.NewKaiClient().ProductClient()).
				Stop(&version.StopOpts{
					ServerName: server,
					ProductID:  productID,
					VersionTag: versionTag,
					Comment:    comment,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
