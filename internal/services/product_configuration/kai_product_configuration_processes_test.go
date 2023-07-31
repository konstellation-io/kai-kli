package productconfiguration_test

import (
	"github.com/konstellation-io/krt/pkg/krt"

	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

func (ch *KaiProductConfigurationTest) TestAddProcess_FirstProcess_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	kaiProductConfig.Workflows[0].Processes = nil
	newProcess := &krt.Process{
		Name:          "new-process",
		Type:          "task",
		Image:         "some/image",
		Replicas:      &_defaultReplicas,
		GPU:           &_defaultGPU,
		Config:        map[string]string{"test1": "value1"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"some-subscription"},
		Networking:    nil,
	}

	// WHEN
	err := kaiProductConfig.AddProcess(_defaultWorkflowName, newProcess)

	// THEN
	ch.Require().NoError(err)
	_, wf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Len(wf.Processes, 1)
	ch.Require().Equal(wf.Processes[0], *newProcess)
}

func (ch *KaiProductConfigurationTest) TestAddProcess_WithExistingProcesses_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newProcess := &krt.Process{
		Name:          "new-process",
		Type:          "task",
		Image:         "some/image",
		Replicas:      &_defaultReplicas,
		GPU:           &_defaultGPU,
		Config:        map[string]string{"test1": "value1"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"some-subscription"},
		Networking:    nil,
	}

	// WHEN
	err := kaiProductConfig.AddProcess(_defaultWorkflowName, newProcess)

	// THEN
	ch.Require().NoError(err)
	_, wf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Len(wf.Processes, 2)
	ch.Require().Equal(wf.Processes[1], *newProcess)
}

func (ch *KaiProductConfigurationTest) TestAddProcess_WorkflowNotExistsExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newProcess := &krt.Process{
		Name:          "new-process",
		Type:          "task",
		Image:         "some/image",
		Replicas:      &_defaultReplicas,
		GPU:           &_defaultGPU,
		Config:        map[string]string{"test1": "value1"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"some-subscription"},
		Networking:    nil,
	}

	// WHEN
	err := kaiProductConfig.AddProcess("some-workflow", newProcess)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
}

func (ch *KaiProductConfigurationTest) TestAddProcess_ProcessAlreadyExists_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newProcess := &krt.Process{
		Name:          _defaultProcessName,
		Type:          "task",
		Image:         "some/image",
		Replicas:      &_defaultReplicas,
		GPU:           &_defaultGPU,
		Config:        map[string]string{"test1": "value1"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"some-subscription"},
		Networking:    nil,
	}

	// WHEN
	err := kaiProductConfig.AddProcess(_defaultWorkflowName, newProcess)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrProcessAlreadyExists)
}

func (ch *KaiProductConfigurationTest) TestUpdateProcess_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newProcess := &krt.Process{
		Name:          _defaultProcessName,
		Type:          "task",
		Image:         "some/image",
		Replicas:      &_defaultReplicas,
		GPU:           &_defaultGPU,
		Config:        map[string]string{"test1": "value1"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"some-subscription"},
		Networking:    nil,
	}

	// WHEN
	err := kaiProductConfig.UpdateProcess(_defaultWorkflowName, newProcess)

	// THEN
	ch.Require().NoError(err)
	_, wf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Len(wf.Processes, 1)
	ch.Require().Equal(wf.Processes[0], *newProcess)
}

func (ch *KaiProductConfigurationTest) TestUpdateProcess_WorkflowNotExistsExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newProcess := &krt.Process{
		Name:          _defaultProcessName,
		Type:          "task",
		Image:         "some/image",
		Replicas:      &_defaultReplicas,
		GPU:           &_defaultGPU,
		Config:        map[string]string{"test1": "value1"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"some-subscription"},
		Networking:    nil,
	}

	// WHEN
	err := kaiProductConfig.UpdateProcess("some-workflow", newProcess)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
}

func (ch *KaiProductConfigurationTest) TestUpdateProcess_ProcessNotExists_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newProcess := &krt.Process{
		Name:          "some-process",
		Type:          "task",
		Image:         "some/image",
		Replicas:      &_defaultReplicas,
		GPU:           &_defaultGPU,
		Config:        map[string]string{"test1": "value1"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"some-subscription"},
		Networking:    nil,
	}

	// WHEN
	err := kaiProductConfig.UpdateProcess(_defaultWorkflowName, newProcess)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrProcessNotFound)
}

func (ch *KaiProductConfigurationTest) TestRemoveProcess_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.RemoveProcess(_defaultWorkflowName, _defaultProcessName)

	// THEN
	ch.Require().NoError(err)
	_, wf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Empty(wf.Processes)
}

func (ch *KaiProductConfigurationTest) TestRemoveProcess_WorkflowNotExistsExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.RemoveProcess("some-workflow", _defaultProcessName)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
}

func (ch *KaiProductConfigurationTest) TestRemoveProcess_ProcessNotExists_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.RemoveProcess(_defaultWorkflowName, "some-process")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrProcessNotFound)
}
