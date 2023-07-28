package configuration

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_serverFlag     = "server"
	_productIDFlag  = "product-id"
	_workflowIDFlag = "workflow-id"
	_processIDFlag  = "process-id"
	_scopeFlag      = "scope"
)

// NewConfigurationCmd creates a new command to handle 'config' subcommands.
func NewConfigurationCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage configuration",
		Example: heredoc.Doc(`
			$ kli config add <key> <value> [--scope <scope>] [--process-id <process_id>] [--workflow-id <workflow_id>] [--product_id <product_id>] [--server <server>]
			$ kli config update <key> <value> [--scope <scope>] [--process-id <process_id>] [--workflow-id <workflow_id>] [--product_id <product_id>] [--server <server>]
			$ kli config remove <key> [--scope <scope>] [--process-id <process_id>] [--workflow-id <workflow_id>] [--product_id <product_id>] [--server <server>]
			$ kli config ls [--scope <scope>] [--process-id <process_id>] [--workflow-id <workflow_id>] [--product_id <product_id>] [--server <server>]
			$ kli config dump [--scope <scope>] [--process-id <process_id>] [--workflow-id <workflow_id>] [--product_id <product_id>] [version_id <version_id>] [--server <server>]
			$ kli config sync [--scope <scope>] [--process-id <process_id>] [--workflow-id <workflow_id>] [--product_id <product_id>] [version_id <version_id>] [--server <server>]
		`),
	}

	cmd.AddCommand(
		NewAddCmd(logger),
		NewRemoveCmd(logger),
		NewListCmd(logger),
	)

	return cmd
}
