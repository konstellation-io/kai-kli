package krt

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/cmd/krt/config"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewKRTCmd creates a new command to handle 'krt' keyword.
func NewKRTCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "krt",
		Short: "Utils for KRT creation and validation",
		Example: heredoc.Doc(`
			$ kli krt create -i/--input <input_path,default=pwd> -o/--output <output_path,default=pwd/<version>.krt>
				[--version <version_name>] [--update-local <true|false, default=true>]
			$ kli krt validate -i/--input <input_path,default=pwd>
		`),
	}

	cmd.AddCommand(
		NewBuildCmd(logger),
		NewValidateCmd(logger),
		config.NewKRTConfigCmd(logger),
	)

	return cmd
}
