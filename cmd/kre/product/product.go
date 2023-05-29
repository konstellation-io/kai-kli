package product

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewProductCmd creates a new command to handle 'product' subcommands.
func NewProductCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "Manage KRE Version",
		Example: heredoc.Doc(`
			$ kli kre product ls
			$ kli kre product create 'product-name' -d 'product description'
		`),
	}

	cmd.PersistentFlags().StringP("server", "s", cfg.DefaultServer, "KRE server to use")

	cmd.AddCommand(
		NewListCmd(logger, cfg),
		NewCreateCmd(logger, cfg),
	)

	return cmd
}
