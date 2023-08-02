package process

import (
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/process"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

// NewUpdateCmd creates a new command to update a process to the given product workflow.
func NewUpdateCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <process_id> [opts...] ",
		Args:    cobra.ExactArgs(1),
		Short:   "Update an existing process for the given product workflow on the given server",
		Example: "$ kli update <process_id> [opts...] ",
		RunE: func(cmd *cobra.Command, args []string) error {
			processID := args[0]

			updateProcessOpts, err := getUpdateProcessOpts(processID, cmd)
			if err != nil {
				return err
			}

			// TODO Get the given product or the default one
			// TODO Get the given server or the default one

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())
			err = process.NewHandler(logger, r).UpdateProcess(updateProcessOpts)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_serverFlag, "", "The server name to update the process.")
	cmd.Flags().String(_productIDFlag, "", "The product ID to update the process.")
	cmd.Flags().String(_workflowIDFlag, "", "The workflow ID to update the process.")
	cmd.Flags().String(_processTypeFlag, "", "The process type to update the process.")
	cmd.Flags().String(_processImageFlag, "", "The process image to update the process.")
	cmd.Flags().Int(_replicasFlag, -1, "The number of replicas to be deployed for the current process.")
	cmd.Flags().String(_cpuRequestFlag, "", "The CPU request resources needed to deploy the current process.")
	cmd.Flags().String(_cpuLimitFlag, "", "The maximum CPU resources allowed to deploy the current process.")
	cmd.Flags().String(_memRequestFlag, "", "The memory request resources needed to deploy the current process.")
	cmd.Flags().String(_memLimitFlag, "", "The maximum memory resources allowed to deploy the current process.")
	cmd.Flags().Int(_networkSourcePort, -1, "The source port used by the current process.")
	cmd.Flags().Int(_networkTargetPort, -1, "The target port exposed by the current process.")
	cmd.Flags().String(_networkProtocol, "", "The network protocol used by the current process.")
	cmd.Flags().StringSlice(_subscriptionsFlag, []string{}, "The subscriptions to be used by the current process.")

	return cmd
}

func getUpdateProcessOpts(processID string, cmd *cobra.Command) (*process.ProcessOpts, error) {
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

	processType, err := cmd.Flags().GetString(_processTypeFlag)
	if err != nil {
		return nil, err
	}

	processImage, err := cmd.Flags().GetString(_processImageFlag)
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
		Image:             processImage,
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
