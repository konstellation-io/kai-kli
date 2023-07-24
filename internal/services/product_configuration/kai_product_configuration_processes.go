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
	wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	if len(wf.Processes) == 0 {
		wf.Processes = []krt.Process{}
	}

	for _, process := range wf.Processes {
		if process.Name == pc.Name {
			return ErrProcessAlreadyExists
		}
	}

	wf.Processes = append(wf.Processes, *pc)

	return nil
}

func (c *KaiProductConfiguration) UpdateProcess(workflowName string, pc *krt.Process) error {
	wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	for i, process := range wf.Processes {
		if process.Name == pc.Name {
			wf.Processes[i] = *pc
			return nil
		}
	}

	return ErrProcessNotFound
}

func (c *KaiProductConfiguration) RemoveProcess(workflowName, processName string) error {
	wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	for i, process := range wf.Processes {
		if process.Name == processName {
			wf.Processes = append(wf.Processes[:i], wf.Processes[i+1:]...)
			return nil
		}
	}

	return ErrProcessNotFound
}
