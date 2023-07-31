package process

import "github.com/konstellation-io/krt/pkg/krt"

type UpdateProcessOpts struct {
	ServerName    string
	ProductID     string
	WorkflowID    string
	ProcessID     string
	ProcessType   krt.ProcessType
	Image         string
	Replicas      int
	GPU           bool
	Subscriptions []string
}

func (w *Handler) UpdateProcess(opts *UpdateProcessOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	existingProcess, err := productConfig.GetProcess(opts.WorkflowID, opts.ProcessID)
	if err != nil {
		return err
	}

	err = productConfig.UpdateProcess(
		opts.WorkflowID,
		&krt.Process{
			Name:          opts.ProcessID,
			Type:          opts.ProcessType,
			Image:         opts.Image,
			Replicas:      &opts.Replicas,
			GPU:           &opts.GPU,
			Config:        existingProcess.Config,
			ObjectStore:   existingProcess.ObjectStore,
			Secrets:       existingProcess.Secrets,
			Subscriptions: opts.Subscriptions,
			Networking:    existingProcess.Networking,
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
