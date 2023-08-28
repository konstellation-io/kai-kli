package processregistry_test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	registry "github.com/konstellation-io/kli/api/process-registry"
	"github.com/konstellation-io/kli/cmd/config"
	processregistry "github.com/konstellation-io/kli/internal/commands/process-registry"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ListProcessSuite struct {
	suite.Suite

	logger             *mocks.MockLogger
	renderer           *mocks.MockRenderer
	processRegistryAPI *mocks.MockProcessRegistryInterface
	manager            *processregistry.Handler
	configuration      *configuration.KaiConfigService
	tmpDir             string
}

func TestListProcessSuite(t *testing.T) {
	suite.Run(t, new(ListProcessSuite))
}

func (s *ListProcessSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	viper.SetDefault(config.KaiProductConfigFolder, ".kai")

	s.renderer = renderer

	s.processRegistryAPI = mocks.NewMockProcessRegistryInterface(ctrl)

	s.manager = processregistry.NewHandler(
		s.logger,
		renderer,
		s.processRegistryAPI,
	)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.tmpDir = tmpDir
	s.configuration = configuration.NewKaiConfigService(s.logger)
}

func (s *ListProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ListProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ListProcessSuite) BeforeTest(_, _ string) {
	err := os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), os.FileMode(0750))
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConf, err := s.configuration.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(&configuration.Server{Name: "test", URL: "http://test.com", IsDefault: true})
	s.Require().NoError(err)

	err = s.configuration.WriteConfiguration(kaiConf)
	s.Require().NoError(err)
}

func (s *ListProcessSuite) TestListProcessRegistries() {
	// GIVEN
	processType := _processType
	productID := _productID
	serverName := _serverName

	testProcessRegistry := registry.ProcessRegistry{
		ID:         "test-process-id",
		Name:       "test-process",
		Version:    "v1.0.0",
		Type:       processType,
		Image:      "test-image",
		UploadDate: time.Now().UTC(),
		Owner:      "test",
	}

	retrievedProcessRegistries := []registry.ProcessRegistry{testProcessRegistry}

	s.processRegistryAPI.EXPECT().List(gomock.Any(), productID, string(processType)).Return(
		retrievedProcessRegistries,
		nil,
	)
	s.renderer.EXPECT().RenderProcessRegistries(retrievedProcessRegistries)

	// WHEN
	err := s.manager.ListProcesses(&processregistry.ListProcessesOpts{
		ServerName:  serverName,
		ProductID:   productID,
		ProcessType: krt.ProcessType(processType),
	})

	// THEN
	s.Require().NoError(err)
}

func (s *ListProcessSuite) TestListProcessRegistriesAPIError() {
	// GIVEN
	processType := _processType
	productID := _productID
	serverName := _serverName

	s.processRegistryAPI.EXPECT().List(gomock.Any(), productID, string(processType)).Return(
		nil,
		fmt.Errorf("mock error"),
	)

	// WHEN
	err := s.manager.ListProcesses(&processregistry.ListProcessesOpts{
		ServerName:  serverName,
		ProductID:   productID,
		ProcessType: krt.ProcessType(processType),
	})

	// THEN
	s.Require().Error(err)
}
