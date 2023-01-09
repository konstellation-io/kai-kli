package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/cmd/factory"
	"github.com/konstellation-io/kli/internal/render"
)

// NewDefaultCmd creates a new command to handle 'default' keyword.
func NewDefaultCmd(f factory.CmdFactory) *cobra.Command {
	log := f.Logger()
	cmd := &cobra.Command{
		Use:   "default <server_name>",
		Short: "Set a default server",
		Long:  "Mark an existing server as default",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := f.Config()
			name := args[0]

			err := cfg.SetDefaultServer(name)
			if err != nil {
				return err
			}

			r := render.DefaultRenderer(cmd.OutOrStdout())
			cfg.RenderServerList(r)

			log.Success(fmt.Sprintf("Server '%s' is now default.\n", name))

			return nil
		},
	}

	return cmd
}
