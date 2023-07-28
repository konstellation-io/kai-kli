package product

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewProductCmd creates a new command to handle 'product' subcommands.
func NewProductCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product <command>",
		Short: "Manage products",
		Example: heredoc.Doc(`
			$ kli products ls [--server <server_name>]
			$ kli product create <product_name> [--server <server_name>] [--version <initial_version>] [--description <description>] [--init-local] [--local-path <local_path>]
			$ kli product bind <product_name> [local_path] [--server <server_name>] [--force]
		`),
	}

	cmd.AddCommand(
		NewCreateCmd(logger),
	)

	return cmd
}
