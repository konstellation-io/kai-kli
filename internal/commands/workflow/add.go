package workflow

import "github.com/konstellation-io/krt/pkg/krt"

type AddWorkflowOpts struct {
	ServerName   string
	ProductID    string
	WorkflowID   string
	WorkflowType krt.WorkflowType
}

func (w *Handler) AddWorkflow(opts *AddWorkflowOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	err = productConfig.AddWorkflow(
		&krt.Workflow{
			Name:      opts.WorkflowID,
			Type:      opts.WorkflowType,
			Config:    map[string]string{},
			Processes: []krt.Process{},
		})

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
