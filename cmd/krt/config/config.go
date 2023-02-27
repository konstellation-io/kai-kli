package config

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

// NewKRTConfigCmd creates a new command to handle 'krt' keyword.
func NewKRTConfigCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Utils for KRT YAML file modification and validation",
		Example: heredoc.Doc(`
			$ kli krt config version demo-v1 -i demo/krt.yaml
			$ kli krt config validate -i demo/krt.yaml
		`),
	}

	cmd.AddCommand(
		NewValidateCmd(logger),
		NewKrtConfigVersionCmd(logger),
	)

	return cmd
}
