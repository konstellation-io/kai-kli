package server_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type AddServerSuite struct {
	suite.Suite

	renderer *mocks.MockRenderer
	manager  *server.Handler
	tmpDir   string
}

func TestAddServerSuite(t *testing.T) {
	suite.Run(t, new(AddServerSuite))
}

func (s *AddServerSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.renderer = renderer
	s.manager = server.NewHandler(logger, renderer)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.tmpDir = tmpDir
}

func (s *AddServerSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *AddServerSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *AddServerSuite) TestAddServer_ValidServerInExistingConfig() {
	existingConfig, err := createDefaultConfiguration()

	s.renderer.EXPECT().RenderServers(gomock.Any()).Times(1)
	s.Require().NoError(err)

	var (
		newServer = server.Server{
			Name:     "valid-server",
			Host:     "valid-kai-server",
			Protocol: server.ProtocolHTTPS,
		}

		expectedConfigServer = configuration.Server{
			Name:       newServer.Name,
			Host:       newServer.Host,
			Protocol:   newServer.Protocol,
			APIURL:     "https://api.valid-kai-server",
			AuthURL:    "https://auth.valid-kai-server",
			StorageURL: "https://storage-console.valid-kai-server",
			IsDefault:  false,
		}

		expectedKaiConfig = configuration.KaiConfiguration{
			Servers: append(existingConfig.Servers, expectedConfigServer),
		}
	)

	err = s.manager.AddNewServer(newServer, false)
	s.Assert().NoError(err)

	configBytes, err := os.ReadFile(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	var kaiConfig configuration.KaiConfiguration

	err = yaml.Unmarshal(configBytes, &kaiConfig)
	s.Require().NoError(err)

	s.Assert().Equal(expectedKaiConfig, kaiConfig)
}

func (s *AddServerSuite) TestAddServer_DefaultServer() {
	existingConfig, err := createDefaultConfiguration()

	s.renderer.EXPECT().RenderServers(gomock.Any()).Times(1)
	s.Require().NoError(err)

	for i := range existingConfig.Servers {
		existingConfig.Servers[i].IsDefault = false
	}

	var (
		newServer = server.Server{
			Name:     "valid-server",
			Host:     "valid-kai-server",
			Protocol: server.ProtocolHTTPS,
		}

		expectedConfigServer = configuration.Server{
			Name:       newServer.Name,
			Host:       newServer.Host,
			Protocol:   server.ProtocolHTTPS,
			APIURL:     "https://api.valid-kai-server",
			AuthURL:    "https://auth.valid-kai-server",
			StorageURL: "https://storage-console.valid-kai-server",
			IsDefault:  true,
		}

		expectedKaiConfig = configuration.KaiConfiguration{
			Servers: append(existingConfig.Servers, expectedConfigServer),
		}
	)

	err = s.manager.AddNewServer(newServer, true)
	s.Assert().NoError(err)

	configBytes, err := os.ReadFile(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	var kaiConfig configuration.KaiConfiguration

	err = yaml.Unmarshal(configBytes, &kaiConfig)
	s.Require().NoError(err)

	s.Assert().Equal(expectedKaiConfig, kaiConfig)
}

func (s *AddServerSuite) TestAddServer_DuplicatedServerName() {
	existingConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	newServer := server.Server{
		Name: existingConfig.Servers[0].Name,
		Host: "new-server.com",
	}

	err = s.manager.AddNewServer(newServer, false)
	s.Assert().ErrorIs(err, configuration.ErrDuplicatedServerName)
}

func (s *AddServerSuite) TestAddServer_DuplicatedServerURL() {
	existingConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	newServer := server.Server{
		Name: "new-server",
		Host: existingConfig.Servers[0].Host,
	}

	err = s.manager.AddNewServer(newServer, false)
	s.Assert().ErrorIs(err, configuration.ErrDuplicatedServerHost)
}

func (s *AddServerSuite) TestAddServer_InvalidHost() {
	newServer := server.Server{
		Name: "new-server",
		Host: "http://invalid-host.com",
	}

	err := s.manager.AddNewServer(newServer, false)
	s.Assert().ErrorIs(err, server.ErrInvalidServerHost)
}

func createDefaultConfiguration() (*configuration.KaiConfiguration, error) {
	defaultConfiguration := configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{
				Name:      "existing-server",
				Host:      "existing-server.com",
				IsDefault: true,
			},
		},
	}

	if err := createConfigWithServer(defaultConfiguration); err != nil {
		return nil, err
	}

	return &defaultConfiguration, nil
}

func createConfigWithServer(cfg configuration.KaiConfiguration) error {
	if err := os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), 0750); err != nil {
		return err
	}

	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(viper.GetString(config.KaiConfigPath), configBytes, 0600)
}
