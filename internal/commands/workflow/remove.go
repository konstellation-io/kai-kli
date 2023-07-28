package workflow

func (w *KaiWorkflow) RemoveWorkflow(serverName, productID, workflowID string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	err = productConfig.RemoveWorkflow(workflowID)

	if err != nil {
		return err
	}

	err = w.configService.WriteConfiguration(productConfig, productID)
	if err != nil {
		return err
	}

	w.renderer.RenderWorkflows(productConfig.Workflows)

	return nil
}
