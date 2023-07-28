package process

import (
	"github.com/konstellation-io/krt/pkg/krt"
)

func (w *Handler) AddProcess(serverName, productID, workflowID, processID, processType, image string,
	replicas int, subscriptions []string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	GPU := false

	err = productConfig.AddProcess(
		workflowID,
		&krt.Process{
			Name:          processID,
			Type:          krt.ProcessType(processType),
			Image:         image,
			Replicas:      &replicas,
			GPU:           &GPU,
			Config:        map[string]string{},
			ObjectStore:   nil,
			Secrets:       map[string]string{},
			Subscriptions: subscriptions,
			Networking:    nil,
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
