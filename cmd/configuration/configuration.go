package configuration

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_serverFlag     = "server"
	_productIDFlag  = "product"
	_workflowIDFlag = "workflow"
	_processIDFlag  = "process"
	_scopeFlag      = "scope"
)

// NewConfigurationCmd creates a new command to handle 'config' subcommands.
func NewConfigurationCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage configuration",
		Example: heredoc.Doc(`
			$ kli config add <key> <value> [opts...]
			$ kli config remove <key> [opts...]
			$ kli config ls [opts...]
		`),
	}

	cmd.AddCommand(
		NewAddCmd(logger),
		NewRemoveCmd(logger),
		NewListCmd(logger),
	)

	return cmd
}
