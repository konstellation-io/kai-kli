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
			for _, process := range workflow.Processes { //nolint:gocritic
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

	if err := c.validateProcess(*i, len(wf.Processes)+1, pc); err != nil {
		return err
	}

	if len(wf.Processes) == 0 {
		c.Workflows[*i].Processes = []krt.Process{}
	}

	for _, process := range wf.Processes { //nolint:gocritic
		if process.Name == pc.Name {
			return ErrProcessAlreadyExists
		}
	}

	if pc.Type != krt.ProcessTypeTrigger {
		pc.Networking = nil
	}

	c.Workflows[*i].Processes = append(c.Workflows[*i].Processes, *pc)

	return nil
}

func (c *KaiProductConfiguration) UpdateProcess(workflowName string, pc *krt.Process) error {
	wi, wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	for i, process := range wf.Processes { //nolint:gocritic
		if process.Name != pc.Name {
			continue
		}

		tmpProcess := c.updateWorkflowInternal(&c.Workflows[*wi].Processes[i], pc)

		if err := c.validateProcess(*wi, i, tmpProcess); err != nil {
			return err
		}

		if tmpProcess.Type != krt.ProcessTypeTrigger {
			tmpProcess.Networking = nil
		}

		c.Workflows[*wi].Processes[i] = *tmpProcess

		return nil
	}

	return ErrProcessNotFound
}

func (c *KaiProductConfiguration) updateWorkflowInternal(current, updated *krt.Process) *krt.Process {
	tmpProcess := current

	if updated.Type != "" {
		tmpProcess.Type = updated.Type
	}

	if updated.Image != "" {
		tmpProcess.Image = updated.Image
	}

	if updated.Replicas != nil {
		tmpProcess.Replicas = updated.Replicas
	}

	if updated.GPU != nil {
		tmpProcess.GPU = updated.GPU
	}

	if updated.ObjectStore != nil {
		tmpProcess.ObjectStore = updated.ObjectStore
	}

	if updated.Subscriptions != nil {
		tmpProcess.Subscriptions = updated.Subscriptions
	}

	if updated.Networking != nil {
		tmpProcess.Networking = updated.Networking
	}

	if updated.ResourceLimits != nil {
		tmpProcess.ResourceLimits = updated.ResourceLimits
	}

	return tmpProcess
}

func (c *KaiProductConfiguration) RemoveProcess(workflowName, processName string) error {
	wi, wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	for i, process := range wf.Processes { //nolint:gocritic
		if process.Name == processName {
			c.Workflows[*wi].Processes = append(c.Workflows[*wi].Processes[:i], c.Workflows[*wi].Processes[i+1:]...)
			return nil
		}
	}

	return ErrProcessNotFound
}

func (c *KaiProductConfiguration) validateProcess(workflowIndex, processIndex int, pc *krt.Process) error {
	return pc.Validate(workflowIndex, processIndex)
}
