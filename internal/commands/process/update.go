package process

import "github.com/konstellation-io/krt/pkg/krt"

func (w *Handler) UpdateProcess(serverName, productID, workflowID, processID, processType, image string,
	replicas int, subscriptions []string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	GPU := false

	existingProcess, err := productConfig.GetProcess(workflowID, processID)
	if err != nil {
		return err
	}

	err = productConfig.UpdateProcess(
		workflowID,
		&krt.Process{
			Name:          processID,
			Type:          krt.ProcessType(processType),
			Image:         image,
			Replicas:      &replicas,
			GPU:           &GPU,
			Config:        existingProcess.Config,
			ObjectStore:   existingProcess.ObjectStore,
			Secrets:       existingProcess.Secrets,
			Subscriptions: subscriptions,
			Networking:    existingProcess.Networking,
		})

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
