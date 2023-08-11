package productconfiguration_test

import (
	"github.com/konstellation-io/krt/pkg/krt"

	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

func (ch *KaiProductConfigurationTest) TestAddWorkflow_FirstWorkflow_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	kaiProductConfig.Workflows = nil
	newWorkflow := &krt.Workflow{
		Name:      "new-process",
		Type:      "data",
		Config:    map[string]string{"test1": "value1"},
		Processes: ch.getDefaultConfiguration().Workflows[0].Processes,
	}

	// WHEN
	err := kaiProductConfig.AddWorkflow(newWorkflow)

	// THEN
	ch.Require().NoError(err)
	ch.Require().Len(kaiProductConfig.Workflows, 1)
	ch.Require().Equal(kaiProductConfig.Workflows[0], *newWorkflow)
}

func (ch *KaiProductConfigurationTest) TestAddWorkflow_WithExistingWorkflow_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newWorkflow := &krt.Workflow{
		Name:      "new-process",
		Type:      "data",
		Config:    map[string]string{"test1": "value1"},
		Processes: ch.getDefaultConfiguration().Workflows[0].Processes,
	}

	// WHEN
	err := kaiProductConfig.AddWorkflow(newWorkflow)

	// THEN
	ch.Require().NoError(err)
	ch.Require().Len(kaiProductConfig.Workflows, 2)
	ch.Require().Equal(kaiProductConfig.Workflows[1], *newWorkflow)
}

func (ch *KaiProductConfigurationTest) TestAddWorkflow_WorkflowAlreadyExists_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newWorkflow := &krt.Workflow{
		Name:      _defaultWorkflowName,
		Type:      "data",
		Config:    map[string]string{"test1": "value1"},
		Processes: ch.getDefaultConfiguration().Workflows[0].Processes,
	}

	// WHEN
	err := kaiProductConfig.AddWorkflow(newWorkflow)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowAlreadyExists)
}

func (ch *KaiProductConfigurationTest) TestUpdateWorkflow_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newWorkflow := &krt.Workflow{
		Name:      _defaultWorkflowName,
		Type:      "data",
		Processes: []krt.Process{},
	}

	// WHEN
	err := kaiProductConfig.UpdateWorkflow(newWorkflow)

	// THEN
	ch.Require().NoError(err)
	ch.Require().Len(kaiProductConfig.Workflows, 1)
	ch.Require().Equal(kaiProductConfig.Workflows[0].Name, newWorkflow.Name)
	ch.Require().Equal(kaiProductConfig.Workflows[0].Type, newWorkflow.Type)
}

func (ch *KaiProductConfigurationTest) TestUpdateWorkflow_WorkflowNotExists_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newWorkflow := &krt.Workflow{
		Name:      "new-workflow",
		Type:      "task",
		Config:    map[string]string{"test1": "value1"},
		Processes: []krt.Process{},
	}

	// WHEN
	err := kaiProductConfig.UpdateWorkflow(newWorkflow)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
}

func (ch *KaiProductConfigurationTest) TestRemoveWorkflow_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.RemoveWorkflow(_defaultWorkflowName)

	// THEN
	ch.Require().NoError(err)
	ch.Require().Empty(kaiProductConfig.Workflows)
}

func (ch *KaiProductConfigurationTest) TestRemoveWorkflow_WorkflowNotExistsExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.RemoveWorkflow("some-workflow")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
}
