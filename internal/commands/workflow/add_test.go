package workflow_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/stretchr/testify/suite"

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

func (s *AddProcessSuite) TestAddWorkflow_ExpectOk() {
	// GIVEN
	server := "server1"
	newWorkflow := krt.Workflow{
		Name:      "My process",
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

/*func (s *AddProcessSuite) TestAddProcess_NonExistingProduct_ExpectError() {
	// GIVEN
	replicas := 1
	gpu := false
	server := "server1"
	product := "product-test"
	workflow := "Workflow-test"
	newProcess := krt.Process{
		Name:          "My process",
		Type:          krt.ProcessTypeTask,
		Image:         "kst/task",
		Replicas:      &replicas,
		GPU:           &gpu,
		Config:        nil,
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"subject1", "subject2"},
		Networking:    nil,
	}

	// WHEN
	err := s.handler.AddProcess(server, product, workflow, newProcess.Name, string(newProcess.Type),
		newProcess.Image, *newProcess.Replicas, newProcess.Subscriptions)

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProductConfigNotFound)
}

func (s *AddProcessSuite) TestAddProcess_NonExistingWorkflow_ExpectError() {
	// GIVEN
	replicas := 1
	gpu := false
	server := "server1"
	workflow := "Workflow-test"
	newProcess := krt.Process{
		Name:          "My process",
		Type:          krt.ProcessTypeTask,
		Image:         "kst/task",
		Replicas:      &replicas,
		GPU:           &gpu,
		Config:        nil,
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"subject1", "subject2"},
		Networking:    nil,
	}

	// WHEN
	err := s.handler.AddProcess(server, s.productName, workflow, newProcess.Name, string(newProcess.Type),
		newProcess.Image, *newProcess.Replicas, newProcess.Subscriptions)

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *AddProcessSuite) TestAddProcess_DuplicatedProcess_ExpectError() {
	// GIVEN
	replicas := 1
	gpu := false
	server := "server1"
	workflow := "Workflow1"
	newProcess := krt.Process{
		Name:          _getDefaultProcess().Name,
		Type:          krt.ProcessTypeTask,
		Image:         "kst/task",
		Replicas:      &replicas,
		GPU:           &gpu,
		Config:        nil,
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"subject1", "subject2"},
		Networking:    nil,
	}

	// WHEN
	err := s.handler.AddProcess(server, s.productName, workflow, newProcess.Name, string(newProcess.Type),
		newProcess.Image, *newProcess.Replicas, newProcess.Subscriptions)

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProcessAlreadyExists)
}*/
