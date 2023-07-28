package process

func (w *Handler) RemoveProcess(serverName, productID, workflowID, processID string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	err = productConfig.RemoveProcess(workflowID, processID)

	if err != nil {
		return err
	}

	err = w.configService.WriteConfiguration(productConfig, productID)
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
