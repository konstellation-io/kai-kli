package server_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/server"
	"github.com/konstellation-io/kli/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type AddServerSuite struct {
	suite.Suite

	manager *server.KaiConfigurator
	tmpDir  string
}

func TestAddServerSuite(t *testing.T) {
	suite.Run(t, new(AddServerSuite))
}

func (s *AddServerSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockServerRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.manager = server.NewKaiConfigurator(logger, renderer)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiPathKey, kaiConfigPath)

	s.tmpDir = tmpDir
}

func (s *AddServerSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *AddServerSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiPathKey))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *AddServerSuite) TestAddServer_InitialConfiguration() {
	var (
		newServer = server.Server{
			Name: "valid-server",
			URL:  "http://valid-kai-server",
		}

		expectedKaiConfig = server.KaiConfiguration{
			DefaultServer: newServer.Name,
			Servers: []server.Server{
				{
					Name: newServer.Name,
					URL:  newServer.URL,
				},
			},
		}
	)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", newServer.URL+"/info",
		httpmock.NewStringResponder(200, `{"isKaiServer": true}`))

	err := s.manager.AddNewServer(newServer, false)
	s.Assert().NoError(err)

	configBytes, err := os.ReadFile(viper.GetString(config.KaiPathKey))
	s.Require().NoError(err)

	var kaiConfig server.KaiConfiguration
	err = yaml.Unmarshal(configBytes, &kaiConfig)
	s.Require().NoError(err)

	s.Assert().Equal(expectedKaiConfig, kaiConfig)
}

func (s *AddServerSuite) TestAddServer_ValidServerInExistingConfig() {
	existingConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	var (
		newServer = server.Server{
			Name: "valid-server",
			URL:  "http://valid-kai-server",
		}

		expectedKaiConfig = server.KaiConfiguration{
			DefaultServer: existingConfig.DefaultServer,
			Servers:       append(existingConfig.Servers, newServer),
		}
	)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", newServer.URL+"/info",
		httpmock.NewStringResponder(200, `{"isKaiServer": true}`))

	err = s.manager.AddNewServer(newServer, false)
	s.Assert().NoError(err)

	configBytes, err := os.ReadFile(viper.GetString(config.KaiPathKey))
	s.Require().NoError(err)

	var kaiConfig server.KaiConfiguration

	err = yaml.Unmarshal(configBytes, &kaiConfig)
	s.Require().NoError(err)

	s.Assert().Equal(expectedKaiConfig, kaiConfig)
}

func (s *AddServerSuite) TestAddServer_DefaultServer() {
	existingConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	var (
		newServer = server.Server{
			Name: "valid-server",
			URL:  "http://valid-kai-server",
		}

		expectedKaiConfig = server.KaiConfiguration{
			DefaultServer: newServer.Name,
			Servers:       append(existingConfig.Servers, newServer),
		}
	)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", newServer.URL+"/info",
		httpmock.NewStringResponder(200, `{"isKaiServer": true}`))

	err = s.manager.AddNewServer(newServer, true)
	s.Assert().NoError(err)

	configBytes, err := os.ReadFile(viper.GetString(config.KaiPathKey))
	s.Require().NoError(err)

	var kaiConfig server.KaiConfiguration

	err = yaml.Unmarshal(configBytes, &kaiConfig)
	s.Require().NoError(err)

	s.Assert().Equal(expectedKaiConfig, kaiConfig)
}

func (s *AddServerSuite) TestAddServer_DuplicatedServerName() {
	existingConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	newServer := server.Server{
		Name: existingConfig.Servers[0].Name,
		URL:  "new-server.com",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", newServer.URL+"/info",
		httpmock.NewStringResponder(200, `{"isKaiServer": true}`))

	err = s.manager.AddNewServer(newServer, false)
	s.Assert().ErrorIs(err, server.ErrDuplicatedServerName)
}

func (s *AddServerSuite) TestAddServer_DuplicatedServerURL() {
	existingConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	newServer := server.Server{
		Name: "new-server",
		URL:  existingConfig.Servers[0].URL,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", newServer.URL+"/info",
		httpmock.NewStringResponder(200, `{"isKaiServer": true}`))

	err = s.manager.AddNewServer(newServer, false)
	s.Assert().ErrorIs(err, server.ErrDuplicatedServerURL)
}

func createDefaultConfiguration() (*server.KaiConfiguration, error) {
	defaultServerName := "existing-server"

	defaultConfiguration := server.KaiConfiguration{
		DefaultServer: defaultServerName,
		Servers: []server.Server{
			{
				Name: defaultServerName,
				URL:  "existing-server.com",
			},
		},
	}

	if err := createConfigWithServer(defaultConfiguration); err != nil {
		return nil, err
	}

	return &defaultConfiguration, nil
}

func createConfigWithServer(cfg server.KaiConfiguration) error {
	if err := os.Mkdir(path.Dir(viper.GetString(config.KaiPathKey)), 0750); err != nil {
		return err
	}

	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(viper.GetString(config.KaiPathKey), configBytes, 0600)
}
