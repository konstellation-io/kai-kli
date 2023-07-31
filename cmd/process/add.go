package process

import (
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/process"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

const (
	_cpuRequestFlag    = "cpu-request"
	_cpuLimitFlag      = "cpu-limit"
	_memRequestFlag    = "mem-request"
	_memLimitFlag      = "mem-limit"
	_replicasFlag      = "replicas"
	_subscriptionsFlag = "subscriptions"
)

// NewAddCmd creates a new command to add a process to the given product workflow.
func NewAddCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <process_id> <process_type> <image> [--server <server_name>] [--product-id <product_id>] [--workflow-id <workflow_id>] [--cpu-request <cpu_request>] [--cpu-limit <cpu_limit>] [--mem-request <mem_request>] [--mem-limit <mem_limit>] [--replicas <replica_number>] [--subscriptions <subscription_id>...] ",
		Args:    cobra.ExactArgs(3), //nolint:gomnd
		Short:   "Registers a new process for the given product workflow on the given server",
		Example: "$ kli process add <process_id> <process_type> <image> [--server <server_name>] [--product-id <product_id>] [--workflow-id <workflow_id>] [--cpu-request <cpu_request>] [--cpu-limit <cpu_limit>] [--mem-request <mem_request>] [--mem-limit <mem_limit>] [--replicas <replica_number>] [--subscriptions <subscription_id>...] ",
		RunE: func(cmd *cobra.Command, args []string) error {
			processID := args[0]
			processType := args[1]
			image := args[2]

			serverName, err := cmd.Flags().GetString(_serverFlag)
			if err != nil {
				serverName = ""
			}

			productID, err := cmd.Flags().GetString(_productIDFlag)
			if err != nil {
				return err
			}

			workflowID, err := cmd.Flags().GetString(_workflowIDFlag)
			if err != nil {
				return err
			}

			replicas, err := cmd.Flags().GetInt(_replicasFlag)
			if err != nil {
				return err
			}

			subscriptions, err := cmd.Flags().GetStringSlice(_subscriptionsFlag)
			if err != nil {
				return err
			}

			// TODO Get the given product or the default one
			// TODO Get the given server or the default one

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = process.NewHandler(logger, r).AddProcess(&process.AddProcessOpts{
				ServerName:    serverName,
				ProductID:     productID,
				WorkflowID:    workflowID,
				ProcessID:     processID,
				ProcessType:   krt.ProcessType(processType),
				Image:         image,
				Replicas:      replicas,
				Subscriptions: subscriptions,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_serverFlag, "", "The server name to register the process.")
	cmd.Flags().String(_productIDFlag, "", "The product ID to register the process.")
	cmd.Flags().String(_workflowIDFlag, "", "The workflow ID to register the process.")
	cmd.Flags().Int(_replicasFlag, 1, "The number of replicas to be deployed for the current process.")
	cmd.Flags().StringSlice(_subscriptionsFlag, []string{}, "The subscriptions to be used by the current process.")

	return cmd
}
