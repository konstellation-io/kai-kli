package kai

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/cmd/kai/product"
	"github.com/konstellation-io/kli/cmd/kai/version"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewKAICmd creates a new command to handle 'kai' keyword.
func NewKAICmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kai",
		Short: "Manage KAI",
		Example: heredoc.Doc(`
			$ kli kai product ls
		`),
	}

	cmd.AddCommand(
		product.NewProductCmd(logger, cfg),
		version.NewVersionCmd(logger, cfg),
	)

	return cmd
}
