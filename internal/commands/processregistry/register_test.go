package processregistry_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/processregistry"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

const (
	_serverName  = "test"
	_processType = "trigger"
	_processID   = "test-process"
	_productID   = "test-product"
	_version     = "v1.0.0"
)

type RegisterProcessSuite struct {
	suite.Suite

	renderer           *mocks.MockRenderer
	logger             *mocks.MockLogger
	manager            *processregistry.Handler
	processRegistryAPI *mocks.MockProcessRegistryInterface
	configuration      *configuration.KaiConfigService
	tmpDir             string
}

func TestRegisterProcessSuite(t *testing.T) {
	suite.Run(t, new(RegisterProcessSuite))
}

func (s *RegisterProcessSuite) SetupSuite() {
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

func (s *RegisterProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *RegisterProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *RegisterProcessSuite) BeforeTest(_, _ string) {
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

func (s *RegisterProcessSuite) TestRegisterNewServer_ValidPaths_ExpectOk() {
	// GIVEN
	registeredProcessID := "process_id"
	serverName := _serverName
	processType := _processType
	processID := _processID
	version := _version
	processPath := "../../../testdata/sample-trigger"
	dockerfilePath := "../../../testdata/sample-trigger/Dockerfile"
	productID := _productID

	s.processRegistryAPI.EXPECT().
		Register(gomock.Any(), gomock.Any(), productID, processID, processType, version).
		Return(registeredProcessID, nil)

	// WHEN
	err := s.manager.RegisterProcess(&processregistry.RegisterProcessOpts{
		ServerName:  serverName,
		ProductID:   productID,
		ProcessType: krt.ProcessType(processType),
		ProcessID:   processID,
		SourcesPath: processPath,
		Dockerfile:  dockerfilePath,
		Version:     version,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *RegisterProcessSuite) TestRegisterNewServer_InvalidSourcePath_ExpectError() {
	// GIVEN
	serverName := _serverName
	processType := _processType
	processID := _processID
	version := _version
	processPath := "../../../testdata/nonexisting-trigger"
	dockerfilePath := "../../../testdata/sample-trigger/Dockerfile"
	productID := _productID

	// WHEN
	err := s.manager.RegisterProcess(&processregistry.RegisterProcessOpts{
		ServerName:  serverName,
		ProductID:   productID,
		ProcessType: krt.ProcessType(processType),
		ProcessID:   processID,
		SourcesPath: processPath,
		Dockerfile:  dockerfilePath,
		Version:     version,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, processregistry.ErrPathDoesNotExist)
}

func (s *RegisterProcessSuite) TestRegisterNewServer_InvalidDockerfilePath_ExpectError() {
	// GIVEN
	serverName := _serverName
	processType := _processType
	processID := _processID
	version := _version
	processPath := "../../../testdata/sample-trigger"
	dockerfilePath := "../../../testdata/nonexisting-trigger/Dockerfile"
	productID := _productID

	// WHEN
	err := s.manager.RegisterProcess(&processregistry.RegisterProcessOpts{
		ServerName:  serverName,
		ProductID:   productID,
		ProcessType: krt.ProcessType(processType),
		ProcessID:   processID,
		SourcesPath: processPath,
		Dockerfile:  dockerfilePath,
		Version:     version,
	})

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, processregistry.ErrPathDoesNotExist)
}
