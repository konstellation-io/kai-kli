package config

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/krt/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	inputFlag  = "input"
	outputFlag = "output"
)

// NewKrtConfigVersionCmd initialize the config version command.
func NewKrtConfigVersionCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version <version_name> [-i/--input input_path] [-o/--output output_path]",
		Args:  cobra.ExactArgs(1),
		Short: "Modify version field on the input krt.yaml",
		Example: heredoc.Doc(`
			$ kli krt config version v1 -i krt.yaml -o krt.yaml
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			versionName := args[0]
			input, err := cmd.Flags().GetString(inputFlag)
			if err != nil {
				return err
			}
			output, err := cmd.Flags().GetString(outputFlag)
			if err != nil {
				return err
			}

			if output == "" {
				output = input
			}

			renderer := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			krtInteractor := config.NewConfigurator(logger, renderer)
			return krtInteractor.UpdateVersion(versionName, input, output)
		},
	}

	pwd, _ := os.Getwd()
	cmd.Flags().StringP(inputFlag, "i", pwd, "Input path to krt.yaml folder")
	cmd.Flags().StringP(outputFlag, "o", "", "Output folder path")

	return cmd
}
