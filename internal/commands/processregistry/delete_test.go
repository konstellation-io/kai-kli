//go:build unit

package processregistry_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/processregistry"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type DeleteProcessSuite struct {
	suite.Suite

	renderer           *mocks.MockRenderer
	logger             *mocks.MockLogger
	manager            *processregistry.Handler
	processRegistryAPI *mocks.MockAPIClient
	configuration      *configuration.KaiConfigService
	tmpDir             string
}

func TestDeleteProcessSuite(t *testing.T) {
	suite.Run(t, new(DeleteProcessSuite))
}

func (s *DeleteProcessSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	viper.SetDefault(config.KaiProductConfigFolder, ".kai")

	s.renderer = renderer

	s.processRegistryAPI = mocks.NewMockAPIClient(ctrl)

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

func (s *DeleteProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *DeleteProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *DeleteProcessSuite) BeforeTest(_, _ string) {
	err := os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), os.FileMode(0750))
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConf, err := s.configuration.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(&configuration.Server{Name: "test", Host: "http://test.com", IsDefault: true})
	s.Require().NoError(err)

	err = s.configuration.WriteConfiguration(kaiConf)
	s.Require().NoError(err)
}

func (s *DeleteProcessSuite) TestDeleteProcess() {
	// GIVEN
	serverName := _serverName
	processID := _processID
	version := _version
	productID := _productID

	deletedProcess := "test-process-id"

	s.processRegistryAPI.EXPECT().Delete(gomock.Any(), productID, processID, version).Return(deletedProcess, nil)
	s.renderer.EXPECT().RenderProcessDeleted(deletedProcess)

	// WHEN
	err := s.manager.DeleteProcess(&processregistry.DeleteProcessOpts{
		ServerName: serverName,
		ProductID:  productID,
		ProcessID:  processID,
		Version:    version,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *DeleteProcessSuite) TestDeletePublicProcess() {
	// GIVEN
	processID := _processID
	version := _version

	deletedProcess := "test-process-id"

	s.processRegistryAPI.EXPECT().DeletePublic(gomock.Any(), processID, version).Return(deletedProcess, nil)
	s.renderer.EXPECT().RenderProcessDeleted(deletedProcess)

	// WHEN
	err := s.manager.DeleteProcess(&processregistry.DeleteProcessOpts{
		ProcessID: processID,
		Version:   version,
		IsPublic:  true,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *DeleteProcessSuite) TestDeleteProcess_InvalidArgs_PublicWithProduct() {
	// GIVEN
	serverName := _serverName
	processID := _processID
	version := _version
	productID := _productID

	// WHEN
	err := s.manager.DeleteProcess(&processregistry.DeleteProcessOpts{
		ServerName: serverName,
		ProductID:  productID,
		ProcessID:  processID,
		Version:    version,
		IsPublic:   true,
	})

	// THEN
	s.Require().Error(err)
}

func (s *DeleteProcessSuite) TestDeleteProcess_InvalidArgs_PrivateWithNoProduct() {
	// GIVEN
	serverName := _serverName
	processID := _processID
	version := _version

	// WHEN
	err := s.manager.DeleteProcess(&processregistry.DeleteProcessOpts{
		ServerName: serverName,
		ProcessID:  processID,
		Version:    version,
	})

	// THEN
	s.Require().Error(err)
}

func (s *DeleteProcessSuite) TestDeleteProcess_ServiceFails() {
	// GIVEN
	serverName := _serverName
	processID := _processID
	version := _version
	productID := _productID

	expectedErr := errors.New("test error")

	s.processRegistryAPI.EXPECT().Delete(gomock.Any(), productID, processID, version).Return("", expectedErr)

	// WHEN
	err := s.manager.DeleteProcess(&processregistry.DeleteProcessOpts{
		ServerName: serverName,
		ProductID:  productID,
		ProcessID:  processID,
		Version:    version,
	})

	// THEN
	s.ErrorIs(err, expectedErr)
}
