package workflow

type ListWorkflowOpts struct {
	ServerName string
	ProductID  string
}

func (w *Handler) ListWorkflows(opts *ListWorkflowOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	w.renderer.RenderWorkflows(productConfig.Workflows)

	return nil
}
