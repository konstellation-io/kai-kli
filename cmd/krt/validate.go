package krt

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/krt/fileutils/compression"
	"github.com/konstellation-io/kli/internal/commands/krt/validate"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

// NewValidateCmd command to create a KRT file from a directory.
func NewValidateCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate -i <input_path>",
		Args:  cobra.ExactArgs(0),
		Short: "Validate a KRT file",
		Example: heredoc.Doc(`
			$ kli krt create -i konstellation/build/version.krt
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			krtFile, err := cmd.Flags().GetString(inputFlag)
			if err != nil {
				return err
			}

			compressor := compression.NewKrtCompressor()
			renderer := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			validator := validate.NewValidator(logger, compressor, renderer)
			err = validator.ValidateKrt(krtFile)
			if err != nil {
				logger.Debug(fmt.Sprintf("Error validating krt: %s", err.Error()))
				return fmt.Errorf("error validating krt file")
			}

			return nil
		},
	}
	cmd.Flags().StringP(inputFlag, "i", "", "Input folder path")
	_ = cmd.MarkFlagRequired(inputFlag)

	return cmd
}
