package krt

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/krt/build"
	"github.com/konstellation-io/kli/internal/commands/krt/fileutils/compression"
	"github.com/konstellation-io/kli/internal/logging"
)

const (
	versionFlag     = "version"
	inputFlag       = "input"
	outputFlag      = "output"
	updateLocalFlag = "update-local"
)

// NewBuildCmd command to create a KRT file from a directory.
func NewBuildCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build -i <input_path> -o <output_path> [-v <version_name>] [-u]",
		Args:  cobra.ExactArgs(0),
		Short: "Build KRT file from a directory",
		Example: heredoc.Doc(`
			$ kli krt create -i konstellation/sample-project -o konstellation/build/version.krt -v v1 --update-local
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := cmd.Flags().GetString(versionFlag)
			if err != nil {
				return err
			}
			src, err := cmd.Flags().GetString(inputFlag)
			if err != nil {
				return err
			}
			target, err := cmd.Flags().GetString(outputFlag)
			if err != nil {
				return err
			}
			updateLocal, err := cmd.Flags().GetBool(updateLocalFlag)
			if err != nil {
				return err
			}

			compressor := compression.NewKrtCompressor()
			builder := build.NewBuilder(logger, compressor)
			err = builder.Build(src, target, version, updateLocal)
			if err != nil {
				return err
			}

			return nil
		},
	}
	pwd, _ := os.Getwd()
	cmd.Flags().StringP(inputFlag, "i", pwd, "Input folder path")
	cmd.Flags().StringP(outputFlag, "o", pwd, "Output path, either a folder or a file")
	cmd.Flags().StringP(versionFlag, "v", "", "KRT version name")
	cmd.Flags().BoolP(updateLocalFlag, "u", false, "Update local krt.yaml")

	return cmd
}
