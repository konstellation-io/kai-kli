package product

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewProductCmd creates a new command to handle 'product' subcommands.
func NewProductCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "Manage KAI Version",
		Example: heredoc.Doc(`
			$ kli kai product ls
			$ kli kai product create 'product-name' -d 'product description'
		`),
	}

	cmd.PersistentFlags().StringP("server", "s", cfg.DefaultServer, "KAI server to use")

	cmd.AddCommand(
		NewListCmd(logger, cfg),
		NewCreateCmd(logger, cfg),
	)

	return cmd
}
