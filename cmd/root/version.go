package root

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// newVersionCmd creates a new command to handle 'version' keyword.
func newVersionCmd(version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprint(cmd.OutOrStdout(), Format(version, buildDate))
		},
	}

	return cmd
}

// Format return the version properly formatted.
func Format(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	if buildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, buildDate)
	}

	return fmt.Sprintf("kli version %s\n", version)
}
