//go:build unit

package productconfiguration_test

import (
	"reflect"
	"testing"

	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/stretchr/testify/suite"

	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

const (
	_defaultProductVersion     = "v1.0.0"
	_defaultProductDescription = "KRT Product description"
	_defaultWorkflowName       = "workflow-1"
	_defaultWorkflowType       = krt.WorkflowTypeData
	_defaultTriggerProcess     = "trigger-1"
	_defaultExitProcess        = "exit-1"
	_defaultProcessType        = krt.ProcessTypeTrigger
	_defaultProcessImage       = "test/trigger-image"
)

//nolint:gochecknoglobals
var (
	_defaultReplicas      = 1
	_defaultGPU           = false
	_defaultConfiguration = map[string]string{"test1": "value1"}
)

type KaiProductConfigurationTest struct {
	suite.Suite
}

func TestKaiProductConfigurationSuite(t *testing.T) {
	suite.Run(t, new(KaiProductConfigurationTest))
}

func (ch *KaiProductConfigurationTest) TestGetVersionConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	versionConfig := kaiProductConfig.GetVersionConfiguration()

	// THEN
	ch.Require().True(reflect.DeepEqual(_defaultConfiguration, versionConfig))
}

func (ch *KaiProductConfigurationTest) TestUpdateVersionConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test2", Value: "value2"}
	expectedConfig := _defaultConfiguration
	expectedConfig[newConfig.Key] = newConfig.Value

	// WHEN
	versionConfig := kaiProductConfig.UpdateVersionConfig(newConfig)

	// THEN
	ch.Require().Len(versionConfig, 2)
	ch.Require().True(reflect.DeepEqual(expectedConfig, versionConfig))
}

func (ch *KaiProductConfigurationTest) TestUpdateVersionConfiguration_UpdateExistingConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test1", Value: "value2"}

	// WHEN
	versionConfig := kaiProductConfig.UpdateVersionConfig(newConfig)

	// THEN
	ch.Require().Len(versionConfig, 1)
	ch.Require().True(reflect.DeepEqual(map[string]string{"test1": "value2"}, versionConfig))
}

func (ch *KaiProductConfigurationTest) TestDeleteVersionConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	kaiProductConfig.UpdateVersionConfig(productconfiguration.ConfigProperty{Key: "test2", Value: "value2"})
	ch.Require().Len(kaiProductConfig.Config, 2)

	// WHEN
	versionConfig := kaiProductConfig.DeleteVersionConfig("test2")

	// THEN
	ch.Require().True(reflect.DeepEqual(_defaultConfiguration, versionConfig))
}

func (ch *KaiProductConfigurationTest) TestDeleteAllVersionConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	kaiProductConfig.UpdateVersionConfig(productconfiguration.ConfigProperty{Key: "test2", Value: "value2"})
	ch.Require().Len(kaiProductConfig.Config, 2)

	// WHEN
	versionConfig := kaiProductConfig.DeleteVersionConfig("test1", "test2")

	// THEN
	ch.Require().Empty(versionConfig)
}

func (ch *KaiProductConfigurationTest) TestGetWorkflowConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	workflowConfig, err := kaiProductConfig.GetWorkflowConfiguration(_defaultWorkflowName)

	// THEN
	ch.Require().NoError(err)
	ch.Require().True(reflect.DeepEqual(_defaultConfiguration, workflowConfig))
}

func (ch *KaiProductConfigurationTest) TestGetWorkflowConfiguration_NonExistentWorkflow_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	workflowConfig, err := kaiProductConfig.GetWorkflowConfiguration("some-non-existent-workflow")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
	ch.Require().Nil(workflowConfig)
}

func (ch *KaiProductConfigurationTest) TestUpdateWorkflowConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test2", Value: "value2"}
	expectedConfig := _defaultConfiguration
	expectedConfig[newConfig.Key] = newConfig.Value

	// WHEN
	workflowConfig, err := kaiProductConfig.UpdateWorkflowConfig(_defaultWorkflowName, newConfig)

	// THEN
	ch.Require().NoError(err)
	_, wConf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Len(wConf.Config, 2)
	ch.Require().True(reflect.DeepEqual(expectedConfig, workflowConfig))
}

func (ch *KaiProductConfigurationTest) TestUpdateWorkflowConfiguration_NonExistingWorkflow_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test2", Value: "value2"}
	expectedConfig := _defaultConfiguration
	expectedConfig[newConfig.Key] = newConfig.Value

	// WHEN
	workflowConfig, err := kaiProductConfig.UpdateWorkflowConfig("some-workflow", newConfig)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
	ch.Require().Nil(workflowConfig)
}

func (ch *KaiProductConfigurationTest) TestUpdateWorkflowConfiguration_UpdateExistingConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test1", Value: "value2"}

	// WHEN
	workflowConfig, err := kaiProductConfig.UpdateWorkflowConfig(_defaultWorkflowName, newConfig)

	// THEN
	ch.Require().NoError(err)
	_, wConf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Len(wConf.Config, 1)
	ch.Require().True(reflect.DeepEqual(map[string]string{"test1": "value2"}, workflowConfig))
}

func (ch *KaiProductConfigurationTest) TestDeleteWorkflowConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	_, err := kaiProductConfig.UpdateWorkflowConfig(_defaultWorkflowName,
		productconfiguration.ConfigProperty{Key: "test2", Value: "value2"})
	ch.Require().NoError(err)

	i, wConf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Equal(0, *i)
	ch.Require().Len(wConf.Config, 2)

	// WHEN
	versionConfig, err := kaiProductConfig.DeleteWorkflowConfig(_defaultWorkflowName, "test2")

	// THEN
	ch.Require().NoError(err)
	i, wConf, err = kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Len(wConf.Config, 1)
	ch.Require().Equal(0, *i)
	ch.Require().True(reflect.DeepEqual(_defaultConfiguration, versionConfig))
}

func (ch *KaiProductConfigurationTest) TestDeleteWorkflowConfiguration_NonExistingWorkflow_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	_, err := kaiProductConfig.UpdateWorkflowConfig(_defaultWorkflowName,
		productconfiguration.ConfigProperty{Key: "test2", Value: "value2"})
	ch.Require().NoError(err)
	i, wConf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Equal(0, *i)
	ch.Require().Len(wConf.Config, 2)

	// WHEN
	versionConfig, err := kaiProductConfig.DeleteWorkflowConfig("some-non-existing-workflow", "test2")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
	ch.Require().Empty(versionConfig)
}

func (ch *KaiProductConfigurationTest) TestDeleteAllWorkflowConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	_, err := kaiProductConfig.UpdateWorkflowConfig(_defaultWorkflowName,
		productconfiguration.ConfigProperty{Key: "test2", Value: "value2"})
	ch.Require().NoError(err)
	i, wConf, err := kaiProductConfig.GetWorkflow(_defaultWorkflowName)
	ch.Require().NoError(err)
	ch.Require().Equal(0, *i)
	ch.Require().Len(wConf.Config, 2)

	// WHEN
	versionConfig, err := kaiProductConfig.DeleteWorkflowConfig(_defaultWorkflowName, "test1", "test2")

	// THEN
	ch.Require().NoError(err)
	ch.Require().Empty(versionConfig)
}

func (ch *KaiProductConfigurationTest) TestGetProcessConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	processConfig, err := kaiProductConfig.GetProcessConfiguration(_defaultWorkflowName, _defaultTriggerProcess)

	// THEN
	ch.Require().NoError(err)
	ch.Require().True(reflect.DeepEqual(_defaultConfiguration, processConfig))
}

func (ch *KaiProductConfigurationTest) TestGetProcessConfiguration_NonExistentWorkflow_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	processConfig, err := kaiProductConfig.GetProcessConfiguration("some-non-existent-workflow", _defaultTriggerProcess)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
	ch.Require().Nil(processConfig)
}

func (ch *KaiProductConfigurationTest) TestGetProcessConfiguration_NonExistentProcess_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	processConfig, err := kaiProductConfig.GetProcessConfiguration(_defaultWorkflowName, "some-non-existent-process")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrProcessNotFound)
	ch.Require().Nil(processConfig)
}

func (ch *KaiProductConfigurationTest) TestUpdateProcessConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test2", Value: "value2"}
	expectedConfig := _defaultConfiguration
	expectedConfig[newConfig.Key] = newConfig.Value

	// WHEN
	workflowConfig, err := kaiProductConfig.UpdateProcessConfig(_defaultWorkflowName, _defaultTriggerProcess, newConfig)

	// THEN
	ch.Require().NoError(err)
	pConf, err := kaiProductConfig.GetProcess(_defaultWorkflowName, _defaultTriggerProcess)
	ch.Require().NoError(err)
	ch.Require().Len(pConf.Config, 2)
	ch.Require().True(reflect.DeepEqual(expectedConfig, workflowConfig))
}

func (ch *KaiProductConfigurationTest) TestUpdateProcessConfiguration_NonExistentWorkflow_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test2", Value: "value2"}

	// WHEN
	workflowConfig, err := kaiProductConfig.
		UpdateProcessConfig("some-nonexisting-workflow", _defaultTriggerProcess, newConfig)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
	ch.Require().Nil(workflowConfig)
}

func (ch *KaiProductConfigurationTest) TestUpdateProcessConfiguration_NonExistentProcess_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()
	newConfig := productconfiguration.ConfigProperty{Key: "test2", Value: "value2"}

	// WHEN
	workflowConfig, err := kaiProductConfig.
		UpdateProcessConfig(_defaultWorkflowName, "some-nonexisting-process", newConfig)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrProcessNotFound)
	ch.Require().Nil(workflowConfig)
}

func (ch *KaiProductConfigurationTest) TestDeleteProcessConfiguration_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	workflowConfig, err := kaiProductConfig.DeleteProcessConfig(_defaultWorkflowName, _defaultTriggerProcess, "test1")

	// THEN
	ch.Require().NoError(err)
	pConf, err := kaiProductConfig.GetProcess(_defaultWorkflowName, _defaultTriggerProcess)
	ch.Require().NoError(err)
	ch.Require().Len(pConf.Config, 0)
	ch.Require().Empty(workflowConfig)
}

func (ch *KaiProductConfigurationTest) TestDeleteProcessConfiguration_NonExistentWorkflow_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	workflowConfig, err := kaiProductConfig.
		DeleteProcessConfig("some-nonexisting-workflow", _defaultTriggerProcess, "test1")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrWorkflowNotFound)
	ch.Require().Nil(workflowConfig)
}

func (ch *KaiProductConfigurationTest) TestDeleteProcessConfiguration_NonExistentProcess_ExpectError() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	workflowConfig, err := kaiProductConfig.
		DeleteProcessConfig(_defaultWorkflowName, "some-nonexisting-process", "test1")

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrProcessNotFound)
	ch.Require().Nil(workflowConfig)
}

func (ch *KaiProductConfigurationTest) getDefaultConfiguration() productconfiguration.KaiProductConfiguration {
	return productconfiguration.KaiProductConfiguration{
		Krt: &krt.Krt{
			Version:     _defaultProductVersion,
			Description: _defaultProductDescription,
			Config:      map[string]string{"test1": "value1"},
			Workflows: []krt.Workflow{
				{
					Name:   _defaultWorkflowName,
					Type:   _defaultWorkflowType,
					Config: map[string]string{"test1": "value1"},
					Processes: []krt.Process{
						{
							Name:          _defaultTriggerProcess,
							Type:          _defaultProcessType,
							Image:         _defaultProcessImage,
							Replicas:      &_defaultReplicas,
							GPU:           &_defaultGPU,
							Config:        map[string]string{"test1": "value1"},
							ObjectStore:   nil,
							Secrets:       nil,
							Subscriptions: []string{},
							ResourceLimits: &krt.ProcessResourceLimits{
								CPU: &krt.ResourceLimit{
									Request: "100m",
									Limit:   "200m",
								},
								Memory: &krt.ResourceLimit{
									Request: "100Mi",
									Limit:   "200Mi",
								},
							},
							Networking: &krt.ProcessNetworking{
								TargetPort:      10000,
								DestinationPort: 15000,
								Protocol:        krt.NetworkingProtocolHTTP,
							},
						},
						{
							Name:          _defaultExitProcess,
							Type:          krt.ProcessTypeExit,
							Image:         _defaultProcessImage,
							Replicas:      &_defaultReplicas,
							GPU:           &_defaultGPU,
							Config:        map[string]string{"test1": "value1"},
							ObjectStore:   nil,
							Secrets:       nil,
							Subscriptions: []string{_defaultTriggerProcess},
							ResourceLimits: &krt.ProcessResourceLimits{
								CPU: &krt.ResourceLimit{
									Request: "100m",
									Limit:   "200m",
								},
								Memory: &krt.ResourceLimit{
									Request: "100Mi",
									Limit:   "200Mi",
								},
							},
						},
					},
				},
			},
		},
	}
}
