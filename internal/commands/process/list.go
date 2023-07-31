package process

type ListProcessOpts struct {
	ServerName string
	ProductID  string
	WorkflowID string
}

func (w *Handler) ListProcesses(opts *ListProcessOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
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
