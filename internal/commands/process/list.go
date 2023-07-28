package process

func (w *Handler) ListProcesses(serverName, productID, workflowID string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	_, wf, err := productConfig.GetWorkflow(workflowID)
	if err != nil {
		return err
	}

	w.renderer.RenderProcesses(wf.Processes)

	return nil
}
