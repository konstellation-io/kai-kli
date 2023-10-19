package version

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

// NewGetCmd creates a new command to get a version of a product given its tag.
func NewGetCmd(logger logging.Interface) *cobra.Command {
	nArgs := 2
	cmd := &cobra.Command{
		Use: "get <product_id> <version_tag>  [opts...]",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Args:  cobra.ExactArgs(nArgs),
		Short: "Get version of a product given a tag",
		Example: `
    	$ kli product version get <product_id> <version_tag> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]
			versionTag := args[1]

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = version.NewHandler(logger, r, api.NewKaiClient().VersionClient()).
				GetVersion(&version.GetVersionOpts{
					ServerName: server,
					ProductID:  productID,
					VersionTag: versionTag,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
