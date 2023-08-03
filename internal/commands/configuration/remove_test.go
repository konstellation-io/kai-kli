package productconfiguration_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	productconfiguration "github.com/konstellation-io/kli/internal/commands/configuration"
	configservice "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
)

type RemoveConfigurationSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *productconfiguration.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestRemoveConfigurationSuite(t *testing.T) {
	suite.Run(t, new(RemoveConfigurationSuite))
}

func (s *RemoveConfigurationSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	viper.SetDefault(config.KaiProductConfigFolder, ".kai")

	s.productName = "my-product"

	s.renderer = renderer

	s.handler = productconfiguration.NewHandler(
		s.logger,
		s.renderer,
	)
}

func (s *RemoveConfigurationSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *RemoveConfigurationSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *RemoveConfigurationSuite) BeforeTest(_, _ string) {
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

func (s *RemoveConfigurationSuite) TestRemoveConfiguration_VersionScope_ExpectOk() {
	// GIVEN
	key := "test1"
	workflow := ""
	process := ""
	scope := "version"
	s.renderer.EXPECT().RenderConfiguration(scope, map[string]string{})

	// WHEN
	err := s.handler.RemoveConfiguration(&productconfiguration.RemoveConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
		Key:        key,
	})

	// THEN
	s.Require().NoError(err)
	config, err := s.productConfig.GetConfiguration(s.productName)
	s.Require().NoError(err)
	s.Assert().Empty(config.GetVersionConfiguration())
}

func (s *RemoveConfigurationSuite) TestRemoveConfiguration_WorkflowScope_ExpectOk() {
	// GIVEN
	key := "test2"
	workflow := "Workflow1"
	process := ""
	scope := "workflow"
	s.renderer.EXPECT().RenderConfiguration(scope, map[string]string{})

	// WHEN
	err := s.handler.RemoveConfiguration(&productconfiguration.RemoveConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
		Key:        key,
	})

	// THEN
	s.Require().NoError(err)
	config, err := s.productConfig.GetConfiguration(s.productName)
	s.Require().NoError(err)
	workflowConfig, err := config.GetWorkflowConfiguration(workflow)
	s.Assert().NoError(err)
	s.Assert().Empty(workflowConfig)
}

func (s *RemoveConfigurationSuite) TestRemoveConfiguration_WorkflowScope_NoWorkflow_ExpectError() {
	// GIVEN
	key := "test2"
	workflow := ""
	process := ""
	scope := "workflow"

	// WHEN
	err := s.handler.RemoveConfiguration(&productconfiguration.RemoveConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
		Key:        key,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *RemoveConfigurationSuite) TestRemoveConfiguration_ProcessScope_ExpectOk() {
	// GIVEN
	key := "test3"
	workflow := "Workflow1"
	process := "Process1"
	scope := "process"
	s.renderer.EXPECT().RenderConfiguration(scope, map[string]string{})

	// WHEN
	err := s.handler.RemoveConfiguration(&productconfiguration.RemoveConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
		Key:        key,
	})

	// THEN
	s.Require().NoError(err)
	config, err := s.productConfig.GetConfiguration(s.productName)
	s.Require().NoError(err)
	processConfig, err := config.GetProcessConfiguration(workflow, process)
	s.Assert().NoError(err)
	s.Assert().Empty(processConfig)
}

func (s *RemoveConfigurationSuite) TestRemoveConfiguration_ProcessScope_NoWorkflow_ExpectError() {
	// GIVEN
	key := "test3"
	workflow := ""
	process := ""
	scope := "process"

	// WHEN
	err := s.handler.RemoveConfiguration(&productconfiguration.RemoveConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
		Key:        key,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *RemoveConfigurationSuite) TestRemoveConfiguration_ProcessScope_NoProcess_ExpectError() {
	// GIVEN
	key := "test3"
	workflow := "Workflow1"
	process := ""
	scope := "process"

	// WHEN
	err := s.handler.RemoveConfiguration(&productconfiguration.RemoveConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
		Key:        key,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProcessNotFound)
}
