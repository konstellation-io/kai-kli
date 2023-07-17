package config

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/krt/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

// NewValidateCmd validates a KRT yaml file.
func NewValidateCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Args:  cobra.ExactArgs(1),
		Short: "Validate a KRT yml file",
		Example: heredoc.Doc(`
			$ kli krt validate -i/--input <input_path,default=pwd>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			krtFile := args[0]

			renderer := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			krtInteractor := config.NewConfigurator(logger, renderer)
			return krtInteractor.Validate(krtFile)
		},
	}

	return cmd
}
