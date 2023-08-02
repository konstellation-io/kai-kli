package process

import (
	"github.com/konstellation-io/krt/pkg/krt"
)

const (
	_defaultCPURequest    = "500m"
	_defaultCPULimit      = "1.5"
	_defaultMemoryRequest = "64Mi"
	_defaultMemoryLimit   = "128Mi"
)

func (w *Handler) AddProcess(opts *ProcessOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	obj := w.getObjectStore(opts)
	network := w.getNetwork(opts)
	CPULimits := w.getCPULimits(opts)
	memoryLimits := w.getMemoryLimits(opts)

	limits := &krt.ProcessResourceLimits{}

	if CPULimits != nil {
		limits.CPU = CPULimits
	}

	if memoryLimits != nil {
		limits.Memory = memoryLimits
	}

	err = productConfig.AddProcess(
		opts.WorkflowID,
		&krt.Process{
			Name:           opts.ProcessID,
			Type:           opts.ProcessType,
			Image:          opts.Image,
			Replicas:       &opts.Replicas,
			GPU:            &opts.GPU,
			Config:         map[string]string{},
			ObjectStore:    obj,
			Secrets:        map[string]string{},
			Subscriptions:  opts.Subscriptions,
			ResourceLimits: limits,
			Networking:     network,
		})

	if err != nil {
		return err
	}

	err = w.configService.WriteConfiguration(productConfig, opts.ProductID)
	if err != nil {
		return err
	}

	_, wf, err := productConfig.GetWorkflow(opts.WorkflowID)
	if err != nil {
		return err
	}

	w.renderer.RenderProcesses(wf.Processes)

	return nil
}
