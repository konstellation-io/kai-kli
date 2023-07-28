package process

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_serverFlag     = "server"
	_productIDFlag  = "product-id"
	_workflowIDFlag = "workflow-id"
)

// NewProcessCmd creates a new command to handle 'process' subcommands.
func NewProcessCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process <command>",
		Short: "Manage processes",
		Example: heredoc.Doc(`
			$ kli process ls [--server <server_name>] [--product_id <product_id>]  [--workflow-id <workflow_id>]
			$ kli process add <process_id> <process_type> <image> [--server <server_name>] [--product-id <product_id>] [--workflow-id <workflow_id>] [--cpu-request <cpu_request>] [--cpu-limit <cpu_limit>] [--mem-request <mem_request>] [--mem-limit <mem_limit>]
			$ kli process update <process_id> <process_type> [--server <server_name>] [--product-id <product_id>] [--workflow-id <workflow_id>] [--process-type <process_type>] [--cpu-request <cpu_request>] [--cpu-limit <cpu_limit>] [--mem-request <mem_request>] [--mem-limit <mem_limit>]
			$ kli process remove <process_id> [--server <server_name>] [--product-id <product_id>] [--workflow-id <workflow_id>]
		`),
	}

	cmd.AddCommand(
		NewAddCmd(logger),
		NewUpdateCmd(logger),
		NewRemoveCmd(logger),
		NewListCmd(logger),
	)

	return cmd
}
