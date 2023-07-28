package productconfiguration_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/stretchr/testify/suite"

	productconfiguration "github.com/konstellation-io/kli/internal/commands/configuration"
	configservice "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
)

type AddConfigurationSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *productconfiguration.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestAddConfigurationSuite(t *testing.T) {
	suite.Run(t, new(AddConfigurationSuite))
}

func (s *AddConfigurationSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	s.productName = "my-product"

	s.renderer = renderer

	s.handler = productconfiguration.NewHandler(
		s.logger,
		s.renderer,
	)
}

func (s *AddConfigurationSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *AddConfigurationSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *AddConfigurationSuite) BeforeTest(_, _ string) {
	replicas := 1
	gpu := false

	err := s.productConfig.WriteConfiguration(
		&configservice.KaiProductConfiguration{
			Krt: &krt.Krt{
				Version: "v0.0.1",
				Config:  map[string]string{"test1": "value1"},
				Workflows: []krt.Workflow{
					{
						Name:   "Workflow1",
						Type:   "data",
						Config: map[string]string{"test2": "value2"},
						Processes: []krt.Process{
							{
								Name:          "Process1",
								Type:          krt.ProcessTypeTrigger,
								Image:         "kst/trigger",
								Replicas:      &replicas,
								GPU:           &gpu,
								Config:        map[string]string{"test3": "value3"},
								ObjectStore:   nil,
								Secrets:       nil,
								Subscriptions: []string{"subject1", "subject2"},
								Networking: &krt.ProcessNetworking{
									TargetPort:      20000,
									DestinationPort: 21000,
									Protocol:        krt.NetworkingProtocolTCP,
								},
							},
						},
					}},
			},
		},
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *AddConfigurationSuite) TestAddConfiguration_VersionScope_ExpectOk() {
	// GIVEN
	key := "newKey"
	value := "newValue"
	workflow := ""
	process := ""
	scope := "version"
	expectedConfig := map[string]string{"test1": "value1", key: value}
	s.renderer.EXPECT().RenderConfiguration(scope, expectedConfig)

	// WHEN
	err := s.handler.
		AddConfiguration(s.productName, workflow, process, scope, key, value)

	// THEN
	s.Require().NoError(err)
	config, err := s.productConfig.GetConfiguration(s.productName)
	s.Require().NoError(err)
	s.Assert().Equal(config.GetVersionConfiguration()[key], value)
}

func (s *AddConfigurationSuite) TestAddConfiguration_WorkflowScope_ExpectOk() {
	// GIVEN
	key := "newKey"
	value := "newValue"
	workflow := "Workflow1"
	process := ""
	scope := "workflow"
	expectedConfig := map[string]string{"test2": "value2", key: value}
	s.renderer.EXPECT().RenderConfiguration(scope, expectedConfig)

	// WHEN
	err := s.handler.
		AddConfiguration(s.productName, workflow, process, scope, key, value)

	// THEN
	s.Require().NoError(err)
	config, err := s.productConfig.GetConfiguration(s.productName)
	s.Require().NoError(err)
	workflowConfig, err := config.GetWorkflowConfiguration(workflow)
	s.Assert().NoError(err)
	s.Assert().Equal(workflowConfig[key], value)
}

func (s *AddConfigurationSuite) TestAddConfiguration_WorkflowScope_NoWorkflow_ExpectError() {
	// GIVEN
	key := "newKey"
	value := "newValue"
	workflow := ""
	process := ""
	scope := "workflow"

	// WHEN
	err := s.handler.
		AddConfiguration(s.productName, workflow, process, scope, key, value)

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *AddConfigurationSuite) TestAddConfiguration_ProcessScope_ExpectOk() {
	// GIVEN
	key := "newKey"
	value := "newValue"
	workflow := "Workflow1"
	process := "Process1"
	scope := "process"
	expectedConfig := map[string]string{"test3": "value3", key: value}
	s.renderer.EXPECT().RenderConfiguration(scope, expectedConfig)

	// WHEN
	err := s.handler.
		AddConfiguration(s.productName, workflow, process, scope, key, value)

	// THEN
	s.Require().NoError(err)
	config, err := s.productConfig.GetConfiguration(s.productName)
	s.Require().NoError(err)
	processConfig, err := config.GetProcessConfiguration(workflow, process)
	s.Assert().NoError(err)
	s.Assert().Equal(processConfig[key], value)
}

func (s *AddConfigurationSuite) TestAddConfiguration_ProcessScope_NoWorkflow_ExpectError() {
	// GIVEN
	key := "newKey"
	value := "newValue"
	workflow := ""
	process := ""
	scope := "process"

	// WHEN
	err := s.handler.
		AddConfiguration(s.productName, workflow, process, scope, key, value)

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *AddConfigurationSuite) TestAddConfiguration_ProcessScope_NoProcess_ExpectError() {
	// GIVEN
	key := "newKey"
	value := "newValue"
	workflow := "Workflow1"
	process := ""
	scope := "process"

	// WHEN
	err := s.handler.
		AddConfiguration(s.productName, workflow, process, scope, key, value)

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProcessNotFound)
}
