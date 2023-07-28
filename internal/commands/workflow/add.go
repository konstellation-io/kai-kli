package workflow

import "github.com/konstellation-io/krt/pkg/krt"

func (w *KaiWorkflow) AddWorkflow(serverName string, productID string, workflowID string, workflowType string) error {
	productConfig, err := w.configService.GetConfiguration(productID)
	if err != nil {
		return err
	}

	err = productConfig.AddWorkflow(
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
