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

type ListConfigurationSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *productconfiguration.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestListConfigurationSuite(t *testing.T) {
	suite.Run(t, new(ListConfigurationSuite))
}

func (s *ListConfigurationSuite) SetupSuite() {
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

func (s *ListConfigurationSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *ListConfigurationSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ListConfigurationSuite) BeforeTest(_, _ string) {
	err := s.productConfig.WriteConfiguration(
		_getDefaultKaiConfig(),
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *ListConfigurationSuite) TestListConfiguration_VersionScope_ExpectOk() {
	// GIVEN
	workflow := ""
	process := ""
	scope := "version"
	s.renderer.EXPECT().RenderConfiguration(scope, _getDefaultKaiConfig().Config)

	// WHEN
	err := s.handler.ListConfiguration(&productconfiguration.ListConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *ListConfigurationSuite) TestListConfiguration_WorkflowScope_ExpectOk() {
	// GIVEN
	workflow := "Workflow1"
	process := ""
	scope := "workflow"
	s.renderer.EXPECT().RenderConfiguration(scope, _getDefaultWorkflow().Config)

	// WHEN
	err := s.handler.ListConfiguration(&productconfiguration.ListConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *ListConfigurationSuite) TestListConfiguration_WorkflowScope_NoWorkflow_ExpectError() {
	// GIVEN
	workflow := ""
	process := ""
	scope := "workflow"

	// WHEN
	err := s.handler.ListConfiguration(&productconfiguration.ListConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *ListConfigurationSuite) TestListConfiguration_ProcessScope_ExpectOk() {
	// GIVEN
	workflow := "Workflow1"
	process := "Process1"
	scope := "process"
	s.renderer.EXPECT().RenderConfiguration(scope, _getDefaultProcess().Config)

	// WHEN
	err := s.handler.ListConfiguration(&productconfiguration.ListConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *ListConfigurationSuite) TestListConfiguration_ProcessScope_NoWorkflow_ExpectError() {
	// GIVEN
	workflow := ""
	process := ""
	scope := "process"

	// WHEN
	err := s.handler.ListConfiguration(&productconfiguration.ListConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *ListConfigurationSuite) TestListConfiguration_ProcessScope_NoProcess_ExpectError() {
	// GIVEN
	workflow := "Workflow1"
	process := ""
	scope := "process"

	// WHEN
	err := s.handler.ListConfiguration(&productconfiguration.ListConfigurationOpts{
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  process,
		Scope:      scope,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProcessNotFound)
}

func _getDefaultKaiConfig() *configservice.KaiProductConfiguration {
	return &configservice.KaiProductConfiguration{
		Krt: &krt.Krt{
			Version:   "v0.0.1",
			Config:    map[string]string{"test1": "value1"},
			Workflows: []krt.Workflow{_getDefaultWorkflow()},
		},
	}
}

func _getDefaultWorkflow() krt.Workflow {
	return krt.Workflow{
		Name:      "Workflow1",
		Type:      "data",
		Config:    map[string]string{"test2": "value2"},
		Processes: []krt.Process{_getDefaultProcess()},
	}
}

func _getDefaultProcess() krt.Process {
	replicas := 1
	gpu := false

	return krt.Process{
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
			Protocol:        krt.NetworkingProtocolHTTP,
		},
	}
}
