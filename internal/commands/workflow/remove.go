package workflow

type RemoveWorkflowOpts struct {
	ServerName string
	ProductID  string
	WorkflowID string
}

func (w *Handler) RemoveWorkflow(opts *RemoveWorkflowOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	err = productConfig.RemoveWorkflow(opts.WorkflowID)

	if err != nil {
		return err
	}

	err = w.configService.WriteConfiguration(productConfig, opts.ProductID)
	if err != nil {
		return err
	}

	w.renderer.RenderWorkflows(productConfig.Workflows)

	return nil
}
