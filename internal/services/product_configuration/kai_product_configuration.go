package productconfiguration

import (
	"errors"

	"github.com/konstellation-io/krt/pkg/krt"
)

var (
	ErrWorkflowNotFound      = errors.New("workflow not found")
	ErrWorkflowAlreadyExists = errors.New("workflow already exists")
)

type ConfigProperty struct {
	Key   string
	Value string
}

type KaiProductConfiguration struct {
	krt.Krt
}

func (c *KaiProductConfiguration) GetProductVersion() string {
	return c.Version
}

func (c *KaiProductConfiguration) GetProductDescription() string {
	return c.Description
}

func (c *KaiProductConfiguration) GetVersionConfiguration() map[string]string {
	return c.Config
}

func (c *KaiProductConfiguration) UpdateVersionConfig(properties ...ConfigProperty) map[string]string {
	if len(properties) > 0 {
		for _, property := range properties {
			c.Config[property.Key] = property.Value
		}
	}

	return c.Config
}
func (c *KaiProductConfiguration) DeleteVersionConfig(keys ...string) map[string]string {
	if len(keys) > 0 {
		for _, key := range keys {
			delete(c.Config, key)
		}
	}

	return c.Config
}

func (c *KaiProductConfiguration) GetWorkflowConfiguration(workflowName string) (map[string]string, error) {
	workflow, err := c.GetWorkflow(workflowName)
	if err != nil {
		return nil, err
	}

	return workflow.Config, nil
}

func (c *KaiProductConfiguration) UpdateWorkflowConfig(workflowName string, properties ...ConfigProperty) (map[string]string, error) {
	wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return nil, err
	}

	if len(properties) > 0 {
		for _, property := range properties {
			wf.Config[property.Key] = property.Value
		}
	}

	return wf.Config, nil
}

func (c *KaiProductConfiguration) DeleteWorkflowConfig(workflowName string, keys ...string) (map[string]string, error) {
	wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return nil, err
	}

	if len(keys) > 0 {
		for _, key := range keys {
			delete(wf.Config, key)
		}
	}

	return wf.Config, nil
}

func (c *KaiProductConfiguration) GetProcessConfiguration(workflowName, processName string) (map[string]string, error) {
	process, err := c.GetProcess(workflowName, processName)
	if err != nil {
		return nil, err
	}

	return process.Config, nil
}
func (c *KaiProductConfiguration) UpdateProcessConfig(workflowName, processName string,
	properties ...ConfigProperty) (map[string]string, error) {
	process, err := c.GetProcess(workflowName, processName)
	if err != nil {
		return nil, err
	}

	if len(properties) > 0 {
		for _, property := range properties {
			process.Config[property.Key] = property.Value
		}
	}

	return process.Config, nil
}

func (c *KaiProductConfiguration) DeleteProcessConfig(workflowName, processName string,
	keys ...string) (map[string]string, error) {
	process, err := c.GetProcess(workflowName, processName)
	if err != nil {
		return nil, err
	}

	if len(keys) > 0 {
		for _, key := range keys {
			delete(process.Config, key)
		}
	}

	return process.Config, nil
}
