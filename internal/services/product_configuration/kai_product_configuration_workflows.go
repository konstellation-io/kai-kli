package productconfiguration

import (
	"github.com/konstellation-io/krt/pkg/krt"
)

func (c *KaiProductConfiguration) GetWorkflow(workflowName string) (*int, *krt.Workflow, error) {
	for i, workflow := range c.Workflows {
		if workflow.Name == workflowName {
			return &i, &workflow, nil
		}
	}

	return nil, nil, ErrWorkflowNotFound
}

func (c *KaiProductConfiguration) AddWorkflow(wf *krt.Workflow) error {
	if len(c.Workflows) == 0 {
		c.Workflows = []krt.Workflow{}
	}

	for _, workflow := range c.Workflows {
		if workflow.Name == wf.Name {
			return ErrWorkflowAlreadyExists
		}
	}

	c.Workflows = append(c.Workflows, *wf)

	return nil
}

func (c *KaiProductConfiguration) UpdateWorkflow(wf *krt.Workflow) error {
	for i, workflow := range c.Workflows {
		if workflow.Name == wf.Name {
			c.Workflows[i] = *wf
			c.Workflows[i].Config = workflow.Config
			c.Workflows[i].Processes = workflow.Processes

			return nil
		}
	}

	return ErrWorkflowNotFound
}

func (c *KaiProductConfiguration) RemoveWorkflow(workflowName string) error {
	for i, workflow := range c.Workflows {
		if workflow.Name == workflowName {
			c.Workflows = append(c.Workflows[:i], c.Workflows[i+1:]...)
			return nil
		}
	}

	return ErrWorkflowNotFound
}
