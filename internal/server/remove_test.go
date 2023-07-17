package server_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/configuration"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/server"
	"github.com/konstellation-io/kli/internal/setup"
	"github.com/konstellation-io/kli/mocks"
)

type RemoveServerSuite struct {
	suite.Suite

	renderer *mocks.MockRenderer
	logger   logging.Interface
	manager  *server.Handler
	tmpDir   string
}

func TestRemoveServerSuite(t *testing.T) {
	suite.Run(t, new(RemoveServerSuite))
}

func (s *RemoveServerSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.logger = logger
	s.renderer = renderer
	s.manager = server.NewServerHandler(logger, renderer)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	err = os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), os.FileMode(0750))
	s.Require().NoError(err)

	s.tmpDir = tmpDir
}

func (s *RemoveServerSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *RemoveServerSuite) BeforeTest(_, _ string) {
	err := setup.NewKaiSetup(s.logger).CreateConfiguration()
	s.Require().NoError(err)
}

func (s *RemoveServerSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *RemoveServerSuite) TestRemoveServer_RemoveExistingServer() {
	// GIVEN
	s.renderer.EXPECT().RenderServers(gomock.Any()).Times(1)
	currentConfig := configuration.NewKaiConfigService(s.logger)
	givenConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{
				Name:      "server1",
				URL:       "http://server1.com",
				IsDefault: true,
			},
			{
				Name:      "server2",
				URL:       "http://server2.com",
				IsDefault: false,
			},
		},
	}
	err := currentConfig.WriteConfiguration(&givenConfig)
	s.Require().NoError(err)

	expectedConfig := &configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{
				Name:      "server1",
				URL:       "http://server1.com",
				IsDefault: true,
			},
		},
	}

	// WHEN
	err = s.manager.RemoveServer("server2")

	// THEN
	s.Require().NoError(err)
	configValues, err := currentConfig.GetConfiguration()
	s.Require().NoError(err)
	s.Equal(expectedConfig, configValues)
}

func (s *RemoveServerSuite) TestRemoveServer_RemoveNonExistingServer() {
	// GIVEN
	currentConfig := configuration.NewKaiConfigService(s.logger)
	givenConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{
				Name:      "server1",
				URL:       "http://server1.com",
				IsDefault: true,
			},
			{
				Name:      "server2",
				URL:       "http://server2.com",
				IsDefault: false,
			},
		},
	}
	err := currentConfig.WriteConfiguration(&givenConfig)
	s.Require().NoError(err)

	// WHEN
	err = s.manager.RemoveServer("server3")

	// THEN
	s.Require().Error(err)
	s.Require().ErrorIs(err, configuration.ErrServerNotFound)
}

func (s *RemoveServerSuite) TestRemoveServer_RemoveDefaultServer() {
	// GIVEN
	s.renderer.EXPECT().RenderServers(gomock.Any()).Times(1)
	currentConfig := configuration.NewKaiConfigService(s.logger)
	givenConfig := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{
				Name:      "server1",
				URL:       "http://server1.com",
				IsDefault: true,
			},
			{
				Name:      "server2",
				URL:       "http://server2.com",
				IsDefault: false,
			},
		},
	}
	err := currentConfig.WriteConfiguration(&givenConfig)
	s.Require().NoError(err)

	expectedConfig := &configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{
				Name:      "server2",
				URL:       "http://server2.com",
				IsDefault: true,
			},
		},
	}

	// WHEN
	err = s.manager.RemoveServer("server1")

	// THEN
	s.Require().NoError(err)
	configValues, err := currentConfig.GetConfiguration()
	s.Require().NoError(err)
	s.Equal(expectedConfig, configValues)
}
