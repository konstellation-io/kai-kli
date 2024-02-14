package process

import "github.com/konstellation-io/krt/pkg/krt"

func (w *Handler) UpdateProcess(opts *ProcessOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	proc, err := productConfig.GetProcess(opts.WorkflowID, opts.ProcessID)
	if err != nil {
		return err
	}

	obj := w.getObjectStore(opts)
	network := w.getNetwork(opts)
	CPULimits := w.updateCPULimits(opts, proc)
	memoryLimits := w.updateMemoryLimits(opts, proc)

	limits := &krt.ProcessResourceLimits{}

	// TODO if the limits are already set and they are not included in the update command,
	// use the ones from the old version of the process
	if CPULimits != nil {
		limits.CPU = CPULimits
	}

	if memoryLimits != nil {
		limits.Memory = memoryLimits
	}

	err = productConfig.UpdateProcess(
		opts.WorkflowID,
		&krt.Process{
			Name:           opts.ProcessID,
			Type:           opts.ProcessType,
			Image:          opts.Image,
			Replicas:       opts.Replicas,
			GPU:            &opts.GPU,
			ObjectStore:    obj,
			Subscriptions:  opts.Subscriptions,
			ResourceLimits: limits,
			Networking:     network,
			NodeSelectors:  w.updateNodeSelectors(opts, proc),
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

func (w *Handler) updateCPULimits(opts *ProcessOpts, pc *krt.Process) *krt.ResourceLimit {
	res := &krt.ResourceLimit{}

	res.Request = _defaultCPURequest
	res.Limit = _defaultCPULimit

	if pc.ResourceLimits.CPU != nil {
		res.Request = pc.ResourceLimits.CPU.Request
		res.Limit = pc.ResourceLimits.CPU.Limit
	}

	if opts.CPURequest != "" {
		res.Request = opts.CPURequest
	}

	if opts.CPULimit != "" {
		res.Limit = opts.CPULimit
	}

	return res
}

func (w *Handler) updateMemoryLimits(opts *ProcessOpts, pc *krt.Process) *krt.ResourceLimit {
	res := &krt.ResourceLimit{}

	res.Request = _defaultMemoryRequest
	res.Limit = _defaultMemoryLimit

	if pc.ResourceLimits.Memory != nil {
		res.Request = pc.ResourceLimits.Memory.Request
		res.Limit = pc.ResourceLimits.Memory.Limit
	}

	if opts.MemoryRequest != "" {
		res.Request = opts.MemoryRequest
	}

	if opts.MemoryLimit != "" {
		res.Limit = opts.MemoryLimit
	}

	return res
}
