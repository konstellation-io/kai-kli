package process_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/process"
	configservice "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
)

type RemoveProcessSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *process.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestRemoveProcessSuite(t *testing.T) {
	suite.Run(t, new(RemoveProcessSuite))
}

func (s *RemoveProcessSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	viper.SetDefault(config.KaiProductConfigFolder, ".kai")

	s.productName = "my-product"

	s.renderer = renderer

	s.handler = process.NewHandler(
		s.logger,
		s.renderer,
	)
}

func (s *RemoveProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *RemoveProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *RemoveProcessSuite) BeforeTest(_, _ string) {
	err := s.productConfig.WriteConfiguration(
		_getDefaultKaiConfig(),
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *RemoveProcessSuite) TestRemoveProcess_ExpectOk() {
	// GIVEN
	server := "server1"
	workflow := "Workflow1"
	processName := _getDefaultProcess().Name

	s.renderer.EXPECT().RenderProcesses(gomock.Eq([]krt.Process{}))

	// WHEN
	err := s.handler.RemoveProcess(&process.RemoveProcessOpts{
		ServerName: server,
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  processName,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *RemoveProcessSuite) TestRemoveProcess_NonExistingProduct_ExpectError() {
	// GIVEN
	server := "server1"
	product := "some-product"
	workflow := "Workflow1"
	processName := _getDefaultProcess().Name

	// WHEN
	err := s.handler.RemoveProcess(&process.RemoveProcessOpts{
		ServerName: server,
		ProductID:  product,
		WorkflowID: workflow,
		ProcessID:  processName,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProductConfigNotFound)
}

func (s *RemoveProcessSuite) TestRemoveProcess_NonExistingWorkflow_ExpectError() {
	// GIVEN
	server := "server1"
	workflow := "some-workflow"
	processName := _getDefaultProcess().Name

	// WHEN
	err := s.handler.RemoveProcess(&process.RemoveProcessOpts{
		ServerName: server,
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  processName,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *RemoveProcessSuite) TestRemoveProcess_NonExistingProcess_ExpectError() {
	// GIVEN
	server := "server1"
	workflow := "Workflow1"
	processName := "some-process"

	// WHEN
	err := s.handler.RemoveProcess(&process.RemoveProcessOpts{
		ServerName: server,
		ProductID:  s.productName,
		WorkflowID: workflow,
		ProcessID:  processName,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProcessNotFound)
}
