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
			$ kli products ls [opts...]
			$ kli product create <product_name> [opts...]
			$ kli product bind <product_name> [opts...]
		`),
	}

	cmd.AddCommand(
		NewCreateCmd(logger),
	)

	return cmd
}
