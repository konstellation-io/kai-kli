package server

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewServerCmd creates a new command to handle 'server' subcommands.
func NewServerCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command>",
		Short: "Manage servers for kli",
		Example: heredoc.Doc(`
			$ kli server ls
<<<<<<<< HEAD:cmd/server/server.go
			$ kli server add my-server http://api.kai.local
========
			$ kli server add my-server http://api.kai.local TOKEN_12345
			$ kli server default my-server
>>>>>>>> origin/develop:cmd/kai/server/server.go
		`),
	}

	cmd.AddCommand(
		NewListCmd(logger),
		NewDefaultCmd(logger),
		NewAddCmd(logger),
	)

	return cmd
}
