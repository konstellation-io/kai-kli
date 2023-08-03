package process

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
)

const (
	_serverFlag           = "server"
	_productIDFlag        = "product"
	_workflowIDFlag       = "workflow"
	_processTypeFlag      = "type"
	_processImageFlag     = "image"
	_gpuFlag              = "gpu"
	_cpuRequestFlag       = "cpu-request"
	_cpuLimitFlag         = "cpu-limit"
	_memRequestFlag       = "mem-request"
	_memLimitFlag         = "mem-limit"
	_replicasFlag         = "replicas"
	_objectStoreNameFlag  = "object-store"
	_objectStoreScopeFlag = "object-store-scope"
	_subscriptionsFlag    = "subscriptions"
	_networkSourcePort    = "source-port"
	_networkTargetPort    = "target-port"
	_networkProtocol      = "protocol"
)

// NewProcessCmd creates a new command to handle 'process' subcommands.
func NewProcessCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process <command> [opts...]",
		Short: "Manage processes",
		Example: heredoc.Doc(`
			$ kli process ls [opts...]
			$ kli process add <process_id> <process_type> <image> [opts...]
			$ kli process update <process_id> <process_type> [opts...]
			$ kli process remove <process_id> [opts...]
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
