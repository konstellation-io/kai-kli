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

type AddProcessSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *workflow.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestAddProcessSuite(t *testing.T) {
	suite.Run(t, new(AddProcessSuite))
}

func (s *AddProcessSuite) SetupSuite() {
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

func (s *AddProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *AddProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *AddProcessSuite) BeforeTest(_, _ string) {
	err := s.productConfig.WriteConfiguration(
		_getDefaultKaiConfig(),
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *AddProcessSuite) TestListWorkflows_ExpectOk() {
	// GIVEN
	server := "server1"
	productID := s.productName
	s.renderer.EXPECT().RenderWorkflows([]krt.Workflow{_getDefaultWorkflow()})

	// WHEN
	err := s.handler.ListWorkflows(&workflow.ListWorkflowOpts{
		ServerName: server,
		ProductID:  productID,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *AddProcessSuite) TestListWorkflows_ProductConfigNotExists_ExpectError() {
	// GIVEN
	server := "server1"
	productID := "non-existing-product"

	// WHEN
	err := s.handler.ListWorkflows(&workflow.ListWorkflowOpts{
		ServerName: server,
		ProductID:  productID,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProductConfigNotFound)
}
