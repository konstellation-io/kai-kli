package process

import "github.com/konstellation-io/krt/pkg/krt"

func (w *Handler) UpdateProcess(opts *ProcessOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	obj := w.getObjectStore(opts)
	network := w.getNetwork(opts)
	CPULimits := w.getCPULimits(opts)
	memoryLimits := w.getMemoryLimits(opts)

	limits := &krt.ProcessResourceLimits{}

	// TODO if the limits are already set and they are not included in the update command,
	// use the ones from the old version of the process
	if CPULimits != nil {
		limits.CPU = CPULimits
	}

	if memoryLimits != nil {
		limits.Memory = memoryLimits
	}

	repl := &opts.Replicas
	if opts.Replicas == -1 {
		repl = nil
	}

	err = productConfig.UpdateProcess(
		opts.WorkflowID,
		&krt.Process{
			Name:           opts.ProcessID,
			Type:           opts.ProcessType,
			Image:          opts.Image,
			Replicas:       repl,
			GPU:            &opts.GPU,
			ObjectStore:    obj,
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
