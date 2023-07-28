package workflow

import "github.com/konstellation-io/krt/pkg/krt"

func (w *KaiWorkflow) UpdateWorkflow(serverName, productID, workflowID, workflowType string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	err = productConfig.UpdateWorkflow(
		krt.Workflow{
			Name:      workflowID,
			Type:      krt.WorkflowType(workflowType),
			Config:    map[string]string{},
			Processes: []krt.Process{},
		})

	if err != nil {
		return err
	}

	err = w.configService.WriteConfiguration(productConfig, productID)
	if err != nil {
		return err
	}

	return nil
}
