//go:build unit

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

type UpdateWorkflowSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *workflow.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestUpdateWorkflowSuite(t *testing.T) {
	suite.Run(t, new(UpdateWorkflowSuite))
}

func (s *UpdateWorkflowSuite) SetupSuite() {
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

func (s *UpdateWorkflowSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *UpdateWorkflowSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *UpdateWorkflowSuite) BeforeTest(_, _ string) {
	err := s.productConfig.WriteConfiguration(
		_getDefaultKaiConfig(),
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *UpdateWorkflowSuite) TestUpdateWorkflow_UpdateFullWorkflow_ExpectOk() {
	// GIVEN
	server := "server1"
	updatedWorkflow := krt.Workflow{
		Name:      _getDefaultWorkflow().Name,
		Type:      krt.WorkflowTypeTraining,
		Config:    _getDefaultWorkflow().Config,
		Processes: _getDefaultWorkflow().Processes,
	}
	s.renderer.EXPECT().RenderWorkflows([]krt.Workflow{updatedWorkflow})

	// WHEN
	err := s.handler.UpdateWorkflow(&workflow.UpdateWorkflowOpts{
		ServerName:   server,
		ProductID:    s.productName,
		WorkflowID:   updatedWorkflow.Name,
		WorkflowType: updatedWorkflow.Type,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *UpdateWorkflowSuite) TestUpdateWorkflow_NoExistingProduct_ExpectError() {
	// GIVEN
	server := "server1"

	// WHEN
	err := s.handler.UpdateWorkflow(&workflow.UpdateWorkflowOpts{
		ServerName:   server,
		ProductID:    "some-product",
		WorkflowID:   _getDefaultWorkflow().Name,
		WorkflowType: _getDefaultWorkflow().Type,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProductConfigNotFound)
}

func (s *UpdateWorkflowSuite) TestUpdateWorkflow_NoExistingWorkflow_ExpectError() {
	// GIVEN
	server := "server1"
	workflowName := "some-workflow"

	// WHEN
	err := s.handler.UpdateWorkflow(&workflow.UpdateWorkflowOpts{
		ServerName:   server,
		ProductID:    s.productName,
		WorkflowID:   workflowName,
		WorkflowType: _getDefaultWorkflow().Type,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}
