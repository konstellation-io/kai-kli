package process_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/internal/commands/process"
	configservice "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ListProcessSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *process.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestListProcessSuite(t *testing.T) {
	suite.Run(t, new(ListProcessSuite))
}

func (s *ListProcessSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	s.productName = "my-product"

	s.renderer = renderer

	s.handler = process.NewHandler(
		s.logger,
		s.renderer,
	)
}

func (s *ListProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *ListProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ListProcessSuite) BeforeTest(_, _ string) {
	err := s.productConfig.WriteConfiguration(
		_getDefaultKaiConfig(),
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *ListProcessSuite) TestListProcess_ExpectOk() {
	// GIVEN
	server := "server1"
	workflow := "Workflow1"

	s.renderer.EXPECT().RenderProcesses(ProcessMatcher(_getDefaultProcess()))

	// WHEN
	err := s.handler.ListProcesses(&process.ListProcessOpts{
		ServerName: server,
		ProductID:  s.productName,
		WorkflowID: workflow,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *ListProcessSuite) TestListProcess_NonExistingProduct_ExpectError() {
	// GIVEN
	server := "server1"
	product := "product-test"
	workflow := "Workflow-test"

	// WHEN
	err := s.handler.ListProcesses(&process.ListProcessOpts{
		ServerName: server,
		ProductID:  product,
		WorkflowID: workflow,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProductConfigNotFound)
}

func (s *ListProcessSuite) TestListProcess_NonExistingWorkflow_ExpectError() {
	// GIVEN
	server := "server1"
	workflow := "Workflow-test"

	// WHEN
	err := s.handler.ListProcesses(&process.ListProcessOpts{
		ServerName: server,
		ProductID:  s.productName,
		WorkflowID: workflow,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}
