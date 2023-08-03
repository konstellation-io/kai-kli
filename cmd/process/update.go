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
	cmd.Flags().String(_objectStoreNameFlag, "", "The object store name to be used by the current process.")
	cmd.Flags().String(_objectStoreScopeFlag, "", "The object store scope to be used by the current process.")
	cmd.Flags().Int(_networkSourcePort, -1, "The source port used by the current process.")
	cmd.Flags().Int(_networkTargetPort, -1, "The target port exposed by the current process.")
	cmd.Flags().String(_networkProtocol, "", "The network protocol used by the current process.")
	cmd.Flags().StringSlice(_subscriptionsFlag, nil, "The subscriptions to be used by the current process.")

	cmd.MarkFlagRequired(_productIDFlag)
	cmd.MarkFlagRequired(_workflowIDFlag)

	return cmd
}

func getUpdateProcessOpts(processID string, cmd *cobra.Command) (*process.ProcessOpts, error) {
	opts := &process.ProcessOpts{
		ProcessID: processID,
	}

	serverName, err := cmd.Flags().GetString(_serverFlag)
	if err != nil {
		serverName = ""
	}

	opts.ServerName = serverName

	productID, err := cmd.Flags().GetString(_productIDFlag)
	if err != nil {
		return nil, err
	}

	opts.ProductID = productID

	workflowID, err := cmd.Flags().GetString(_workflowIDFlag)
	if err != nil {
		return nil, err
	}

	opts.WorkflowID = workflowID

	if cmd.Flag(_processTypeFlag).Changed {
		processType, err := cmd.Flags().GetString(_processTypeFlag)
		if err != nil {
			return nil, err
		}

		opts.ProcessType = krt.ProcessType(processType)
	}

	if cmd.Flag(_processImageFlag).Changed {
		processImage, err := cmd.Flags().GetString(_processImageFlag)
		if err != nil {
			return nil, err
		}

		opts.Image = processImage
	}

	if cmd.Flag(_processTypeFlag).Changed {
		cpuRequest, err := cmd.Flags().GetString(_cpuRequestFlag)
		if err != nil {
			return nil, err
		}

		opts.CPURequest = cpuRequest
	}

	if cmd.Flag(_processTypeFlag).Changed {
		cpuLimit, err := cmd.Flags().GetString(_cpuLimitFlag)
		if err != nil {
			return nil, err
		}

		opts.CPULimit = cpuLimit
	}

	if cmd.Flag(_memRequestFlag).Changed {
		memRequest, err := cmd.Flags().GetString(_memRequestFlag)
		if err != nil {
			return nil, err
		}

		opts.MemoryRequest = memRequest
	}

	if cmd.Flag(_memLimitFlag).Changed {
		memLimit, err := cmd.Flags().GetString(_memLimitFlag)
		if err != nil {
			return nil, err
		}

		opts.MemoryLimit = memLimit
	}

	if cmd.Flag(_objectStoreNameFlag).Changed {
		objectStoreName, err := cmd.Flags().GetString(_objectStoreNameFlag)
		if err != nil {
			return nil, err
		}

		opts.ObjectStoreName = objectStoreName
	}

	if cmd.Flag(_objectStoreScopeFlag).Changed {
		objectStoreScope, err := cmd.Flags().GetString(_objectStoreScopeFlag)
		if err != nil {
			return nil, err
		}

		opts.ObjectStoreScope = krt.ObjectStoreScope(objectStoreScope)
	}

	if cmd.Flag(_replicasFlag).Changed {
		replicas, err := cmd.Flags().GetInt(_replicasFlag)
		if err != nil {
			return nil, err
		}

		opts.Replicas = &replicas
	}

	if cmd.Flag(_subscriptionsFlag).Changed {
		subscriptions, err := cmd.Flags().GetStringSlice(_subscriptionsFlag)
		if err != nil {
			return nil, err
		}

		opts.Subscriptions = subscriptions
	}

	if cmd.Flag(_networkSourcePort).Changed {
		sourcePort, err := cmd.Flags().GetInt(_networkSourcePort)
		if err != nil {
			return nil, err
		}

		opts.NetworkSourcePort = &sourcePort
	}

	if cmd.Flag(_networkTargetPort).Changed {
		targetPort, err := cmd.Flags().GetInt(_networkTargetPort)
		if err != nil {
			return nil, err
		}

		opts.NetworkTargetPort = &targetPort
	}

	if cmd.Flag(_networkProtocol).Changed {
		protocol, err := cmd.Flags().GetString(_networkProtocol)
		if err != nil {
			return nil, err
		}

		opts.NetworkProtocol = krt.NetworkingProtocol(protocol)
	}

	return opts, nil
}
