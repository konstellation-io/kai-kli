package process

import (
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/process"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

// NewAddCmd creates a new command to add a process to the given product workflow.
func NewAddCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <process_id> <process_type> <image> [opts...] ",
		Args:    cobra.ExactArgs(3), //nolint:gomnd
		Short:   "Registers a new process for the given product workflow on the given server",
		Example: "$ kli process add <process_id> <process_type> <image> [opts...] ",
		RunE: func(cmd *cobra.Command, args []string) error {
			processID := args[0]
			processType := args[1]
			image := args[2]

			addProcessOpts, err := getAddProcessOpts(processID, processType, image, cmd)
			if err != nil {
				return err
			}

			// TODO Get the given product or the default one, Get the given server or the default one

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = process.NewHandler(logger, r).AddProcess(addProcessOpts)
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
	cmd.Flags().String(_cpuRequestFlag, "", "The CPU request resources needed to deploy the current process.")
	cmd.Flags().String(_cpuLimitFlag, "", "The maximum CPU resources allowed to deploy the current process.")
	cmd.Flags().String(_memRequestFlag, "", "The memory request resources needed to deploy the current process.")
	cmd.Flags().String(_memLimitFlag, "", "The maximum memory resources allowed to deploy the current process.")
	cmd.Flags().Int(_networkSourcePort, -1, "The source port used by the current process.")
	cmd.Flags().Int(_networkTargetPort, -1, "The target port exposed by the current process.")
	cmd.Flags().String(_networkProtocol, "", "The network protocol used by the current process.")
	cmd.Flags().StringSlice(_subscriptionsFlag, []string{}, "The subscriptions to be used by the current process.")

	cmd.MarkFlagRequired(_productIDFlag)  //nolint:errcheck
	cmd.MarkFlagRequired(_workflowIDFlag) //nolint:errcheck
	cmd.MarkFlagRequired(_cpuRequestFlag) //nolint:errcheck
	cmd.MarkFlagRequired(_memRequestFlag) //nolint:errcheck

	return cmd
}

func getAddProcessOpts(processID, processType, image string, cmd *cobra.Command) (*process.ProcessOpts, error) {
	serverName, err := cmd.Flags().GetString(_serverFlag)
	if err != nil {
		serverName = ""
	}

	productID, err := cmd.Flags().GetString(_productIDFlag)
	if err != nil {
		return nil, err
	}

	workflowID, err := cmd.Flags().GetString(_workflowIDFlag)
	if err != nil {
		return nil, err
	}

	cpuRequest, err := cmd.Flags().GetString(_cpuRequestFlag)
	if err != nil {
		return nil, err
	}

	cpuLimit, err := cmd.Flags().GetString(_cpuLimitFlag)
	if err != nil {
		return nil, err
	}

	memRequest, err := cmd.Flags().GetString(_memRequestFlag)
	if err != nil {
		return nil, err
	}

	memLimit, err := cmd.Flags().GetString(_memLimitFlag)
	if err != nil {
		return nil, err
	}

	replicas, err := cmd.Flags().GetInt(_replicasFlag)
	if err != nil {
		return nil, err
	}

	subscriptions, err := cmd.Flags().GetStringSlice(_subscriptionsFlag)
	if err != nil {
		return nil, err
	}

	sourcePort, err := cmd.Flags().GetInt(_networkSourcePort)
	if err != nil {
		return nil, err
	}

	targetPort, err := cmd.Flags().GetInt(_networkTargetPort)
	if err != nil {
		return nil, err
	}

	protocol, err := cmd.Flags().GetString(_networkProtocol)
	if err != nil {
		return nil, err
	}

	return &process.ProcessOpts{
		ServerName:        serverName,
		ProductID:         productID,
		WorkflowID:        workflowID,
		ProcessID:         processID,
		ProcessType:       krt.ProcessType(processType),
		Image:             image,
		Replicas:          replicas,
		Subscriptions:     subscriptions,
		CPURequest:        cpuRequest,
		CPULimit:          cpuLimit,
		MemoryRequest:     memRequest,
		MemoryLimit:       memLimit,
		NetworkSourcePort: sourcePort,
		NetworkTargetPort: targetPort,
		NetworkProtocol:   krt.NetworkingProtocol(protocol),
	}, nil
}
