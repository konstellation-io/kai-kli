package workflow

import "github.com/konstellation-io/krt/pkg/krt"

type UpdateWorkflowOpts struct {
	ServerName   string
	ProductID    string
	WorkflowID   string
	WorkflowType krt.WorkflowType
}

func (w *Handler) UpdateWorkflow(opts *UpdateWorkflowOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	err = productConfig.UpdateWorkflow(
		&krt.Workflow{
			Name: opts.WorkflowID,
			Type: opts.WorkflowType,
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
