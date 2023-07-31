package process

type RemoveProcessOpts struct {
	ServerName string
	ProductID  string
	WorkflowID string
	ProcessID  string
}

func (w *Handler) RemoveProcess(opts *RemoveProcessOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	err = productConfig.RemoveProcess(opts.WorkflowID, opts.ProcessID)

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
