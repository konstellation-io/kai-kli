package workflow_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/workflow"
	configservice "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
)

type AddWorkflowSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *workflow.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestAddWorkflowSuite(t *testing.T) {
	suite.Run(t, new(AddWorkflowSuite))
}

func (s *AddWorkflowSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	viper.SetDefault(config.KaiProductConfigFolder, ".kai")

	s.productName = "my-product"

	s.renderer = renderer

	s.handler = workflow.NewHandler(
		s.logger,
		s.renderer,
	)
}

func (s *AddWorkflowSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *AddWorkflowSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *AddWorkflowSuite) BeforeTest(_, _ string) {
	err := s.productConfig.WriteConfiguration(
		_getDefaultKaiConfig(),
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *AddWorkflowSuite) TestAddWorkflow_ExpectOk() {
	// GIVEN
	server := "server1"
	newWorkflow := krt.Workflow{
		Name:      "my-process",
		Type:      krt.WorkflowTypeData,
		Config:    map[string]string{},
		Processes: []krt.Process{},
	}
	s.renderer.EXPECT().RenderWorkflows([]krt.Workflow{_getDefaultWorkflow(), newWorkflow})

	// WHEN
	err := s.handler.AddWorkflow(&workflow.AddWorkflowOpts{
		ServerName:   server,
		ProductID:    s.productName,
		WorkflowID:   newWorkflow.Name,
		WorkflowType: newWorkflow.Type,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *AddWorkflowSuite) TestAddWorkflow_NoExistingProduct_ExpectError() {
	// GIVEN
	server := "server1"
	newWorkflow := krt.Workflow{
		Name:      "my-process",
		Type:      krt.WorkflowTypeData,
		Config:    map[string]string{},
		Processes: []krt.Process{},
	}

	// WHEN
	err := s.handler.AddWorkflow(&workflow.AddWorkflowOpts{
		ServerName:   server,
		ProductID:    "some-product",
		WorkflowID:   newWorkflow.Name,
		WorkflowType: newWorkflow.Type,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProductConfigNotFound)
}

func (s *AddWorkflowSuite) TestAddWorkflow_DuplicatedWorkflow_ExpectError() {
	// GIVEN
	server := "server1"
	newWorkflow := krt.Workflow{
		Name:      _getDefaultWorkflow().Name,
		Type:      krt.WorkflowTypeData,
		Config:    map[string]string{},
		Processes: []krt.Process{},
	}

	// WHEN
	err := s.handler.AddWorkflow(&workflow.AddWorkflowOpts{
		ServerName:   server,
		ProductID:    s.productName,
		WorkflowID:   newWorkflow.Name,
		WorkflowType: newWorkflow.Type,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowAlreadyExists)
}
