package product

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/cmd/product/version"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewProductCmd creates a new command to handle 'product' subcommands.
func NewProductCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product <command>",
		Short: "Manage products",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Example: heredoc.Doc(`
			$ kli products ls [opts...]
			$ kli product create <product_name> [opts...]
			$ kli product bind <product_name> [opts...]
			$ kli product add-user <product_name> <user-email> [opts...]
			$ kli product remove-user <product_name> <user-email> [opts...]
			$ kli product add-maintainer <product_name> <user-email> [opts...]
			$ kli product remove-maintainer <product_name> <user-email> [opts...]
		`),
	}

	cmd.PersistentFlags().StringP("server", "s", "", "KAI server to use")

	cmd.AddCommand(
		NewCreateCmd(logger),
		NewBindCmd(logger),
		NewListCmd(logger),
		NewAddUserCmd(logger),
		NewRemoveUserCmd(logger),
		NewAddMaintainerCmd(logger),
		NewRemoveMaintainerCmd(logger),
		version.NewVersionCmd(logger),
	)

	return cmd
}
