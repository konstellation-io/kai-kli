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

type UpdateProcessSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	logger        *mocks.MockLogger
	handler       *process.Handler
	productConfig *configservice.ProductConfigService
	productName   string
}

func TestUpdateProcessSuite(t *testing.T) {
	suite.Run(t, new(UpdateProcessSuite))
}

func (s *UpdateProcessSuite) SetupSuite() {
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

func (s *UpdateProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)
}

func (s *UpdateProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *UpdateProcessSuite) BeforeTest(_, _ string) {
	err := s.productConfig.WriteConfiguration(
		_getDefaultKaiConfig(),
		s.productName,
	)
	s.Require().NoError(err)
}

func (s *UpdateProcessSuite) TestUpdateProcess_ExpectOk() {
	// GIVEN
	replicas := 2
	server := "server1"
	workflow := "Workflow1"
	newProcess := krt.Process{
		Name:           _getDefaultProcess().Name,
		Type:           krt.ProcessTypeTask,
		Image:          "kst/task2",
		Replicas:       &replicas,
		GPU:            _getDefaultProcess().GPU,
		Config:         _getDefaultProcess().Config,
		ObjectStore:    _getDefaultProcess().ObjectStore,
		Secrets:        _getDefaultProcess().Secrets,
		Subscriptions:  []string{"subject1", "subject2"},
		ResourceLimits: _getDefaultProcess().ResourceLimits,
		Networking:     _getDefaultProcess().Networking,
	}
	s.renderer.EXPECT().RenderProcesses(ProcessMatcher(newProcess))

	// WHEN
	err := s.handler.
		UpdateProcess(&process.ProcessOpts{
			ServerName:        server,
			ProductID:         s.productName,
			WorkflowID:        workflow,
			ProcessID:         newProcess.Name,
			ProcessType:       newProcess.Type,
			Image:             newProcess.Image,
			Replicas:          newProcess.Replicas,
			Subscriptions:     newProcess.Subscriptions,
			NetworkTargetPort: &newProcess.Networking.TargetPort,
			NetworkSourcePort: &newProcess.Networking.DestinationPort,
			NetworkProtocol:   newProcess.Networking.Protocol,
		})

	// THEN
	s.Require().NoError(err)
}

func (s *UpdateProcessSuite) TestUpdateProcess_NonExistingProduct_ExpectError() {
	// GIVEN
	replicas := 1
	server := "server1"
	product := "product-test"
	workflow := "Workflow-test"
	newProcess := krt.Process{
		Name:          _getDefaultProcess().Name,
		Type:          krt.ProcessTypeTask,
		Image:         "kst/task",
		Replicas:      &replicas,
		GPU:           _getDefaultProcess().GPU,
		Config:        _getDefaultProcess().Config,
		ObjectStore:   _getDefaultProcess().ObjectStore,
		Secrets:       _getDefaultProcess().Secrets,
		Subscriptions: []string{"subject1", "subject2"},
		Networking:    _getDefaultProcess().Networking,
	}

	// WHEN
	err := s.handler.UpdateProcess(&process.ProcessOpts{
		ServerName:        server,
		ProductID:         product,
		WorkflowID:        workflow,
		ProcessID:         newProcess.Name,
		ProcessType:       newProcess.Type,
		Image:             newProcess.Image,
		Replicas:          newProcess.Replicas,
		Subscriptions:     newProcess.Subscriptions,
		NetworkTargetPort: &newProcess.Networking.TargetPort,
		NetworkSourcePort: &newProcess.Networking.DestinationPort,
		NetworkProtocol:   newProcess.Networking.Protocol,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProductConfigNotFound)
}

func (s *UpdateProcessSuite) TestUpdateProcess_NonExistingWorkflow_ExpectError() {
	// GIVEN
	replicas := 1
	server := "server1"
	workflow := "Workflow-test"
	newProcess := krt.Process{
		Name:          _getDefaultProcess().Name,
		Type:          krt.ProcessTypeTask,
		Image:         "kst/task",
		Replicas:      &replicas,
		GPU:           _getDefaultProcess().GPU,
		Config:        _getDefaultProcess().Config,
		ObjectStore:   _getDefaultProcess().ObjectStore,
		Secrets:       _getDefaultProcess().Secrets,
		Subscriptions: []string{"subject1", "subject2"},
		Networking:    _getDefaultProcess().Networking,
	}

	// WHEN
	err := s.handler.UpdateProcess(&process.ProcessOpts{
		ServerName:        server,
		ProductID:         s.productName,
		WorkflowID:        workflow,
		ProcessID:         newProcess.Name,
		ProcessType:       newProcess.Type,
		Image:             newProcess.Image,
		Replicas:          newProcess.Replicas,
		Subscriptions:     newProcess.Subscriptions,
		NetworkTargetPort: &newProcess.Networking.TargetPort,
		NetworkSourcePort: &newProcess.Networking.DestinationPort,
		NetworkProtocol:   newProcess.Networking.Protocol,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrWorkflowNotFound)
}

func (s *UpdateProcessSuite) TestUpdateProcess_NonExistingProcess_ExpectError() {
	// GIVEN
	replicas := 1
	server := "server1"
	workflow := "Workflow1"
	newProcess := krt.Process{
		Name:          "My new process",
		Type:          krt.ProcessTypeTask,
		Image:         "kst/task",
		Replicas:      &replicas,
		GPU:           _getDefaultProcess().GPU,
		Config:        _getDefaultProcess().Config,
		ObjectStore:   _getDefaultProcess().ObjectStore,
		Secrets:       _getDefaultProcess().Secrets,
		Subscriptions: []string{"subject1", "subject2"},
		Networking:    _getDefaultProcess().Networking,
	}

	// WHEN
	err := s.handler.UpdateProcess(&process.ProcessOpts{
		ServerName:        server,
		ProductID:         s.productName,
		WorkflowID:        workflow,
		ProcessID:         newProcess.Name,
		ProcessType:       newProcess.Type,
		Image:             newProcess.Image,
		Replicas:          newProcess.Replicas,
		Subscriptions:     newProcess.Subscriptions,
		NetworkTargetPort: &newProcess.Networking.TargetPort,
		NetworkSourcePort: &newProcess.Networking.DestinationPort,
		NetworkProtocol:   newProcess.Networking.Protocol,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configservice.ErrProcessNotFound)
}
