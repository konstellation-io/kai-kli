package workflow

func (w *KaiWorkflow) ListWorkflows(serverName, productID string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	w.renderer.RenderWorkflows(productConfig.Workflows)

	return nil
}
