package productconfiguration

import (
	"errors"

	"github.com/konstellation-io/krt/pkg/krt"
)

var (
	ErrProcessNotFound      = errors.New("process not found")
	ErrProcessAlreadyExists = errors.New("process already exists")
)

func (c *KaiProductConfiguration) GetProcess(workflowName, processName string) (*krt.Process, error) {
	for _, workflow := range c.Workflows {
		if workflow.Name == workflowName {
			for _, process := range workflow.Processes {
				if process.Name == processName {
					return &process, nil
				}
			}

			return nil, ErrProcessNotFound
		}
	}

	return nil, ErrWorkflowNotFound
}

func (c *KaiProductConfiguration) AddProcess(workflowName string, pc *krt.Process) error {
	i, wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	if len(wf.Processes) == 0 {
		c.Workflows[*i].Processes = []krt.Process{}
	}

	for _, process := range wf.Processes {
		if process.Name == pc.Name {
			return ErrProcessAlreadyExists
		}
	}

	c.Workflows[*i].Processes = append(c.Workflows[*i].Processes, *pc)

	return nil
}

func (c *KaiProductConfiguration) UpdateProcess(workflowName string, pc *krt.Process) error {
	wi, wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	for i, process := range wf.Processes {
		if process.Name == pc.Name {
			c.Workflows[*wi].Processes[i] = *pc
			c.Workflows[*wi].Processes[i].Config = pc.Config
			c.Workflows[*wi].Processes[i].GPU = pc.GPU
			c.Workflows[*wi].Processes[i].ObjectStore = pc.ObjectStore
			c.Workflows[*wi].Processes[i].Secrets = pc.Secrets
			c.Workflows[*wi].Processes[i].Subscriptions = pc.Subscriptions
			c.Workflows[*wi].Processes[i].Networking = pc.Networking

			return nil
		}
	}

	return ErrProcessNotFound
}

func (c *KaiProductConfiguration) RemoveProcess(workflowName, processName string) error {
	wi, wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	for i, process := range wf.Processes {
		if process.Name == processName {
			c.Workflows[*wi].Processes = append(c.Workflows[*wi].Processes[:i], c.Workflows[*wi].Processes[i+1:]...)
			return nil
		}
	}

	return ErrProcessNotFound
}
