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
	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type SetDefaultServerSuite struct {
	suite.Suite

	renderer       *mocks.MockRenderer
	manager        *server.Handler
	logger         logging.Interface
	tmpDir         string
	kaiConfService *configuration.KaiConfigService
}

func TestSetDefaultServerSuite(t *testing.T) {
	suite.Run(t, new(SetDefaultServerSuite))
}

func (s *SetDefaultServerSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.logger = logger
	s.renderer = renderer
	s.manager = server.NewHandler(logger, renderer)

	tmpDir, err := os.MkdirTemp("", "TestSetDefaultServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.kaiConfService = configuration.NewKaiConfigService(s.logger)

	s.tmpDir = tmpDir
}

func (s *SetDefaultServerSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *SetDefaultServerSuite) BeforeTest(_, _ string) {
	err := os.MkdirAll(path.Dir(viper.GetString(config.KaiConfigPath)), 0750)
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConf, err := s.kaiConfService.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(&configuration.Server{
		Name:      "my-server1",
		URL:       "https://kai-dev.konstellation.io",
		AuthURL:   "https://auth.kai-dev.konstellation1.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
		IsDefault: false,
		Token:     &configuration.Token{},
	})
	s.Require().NoError(err)

	err = kaiConf.AddServer(&configuration.Server{
		Name:      "my-server2",
		URL:       "https://kai-dev.konstellation2.io",
		AuthURL:   "https://auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
		IsDefault: true,
		Token:     &configuration.Token{},
	})
	s.Require().NoError(err)

	err = s.kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)
}

func (s *SetDefaultServerSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *SetDefaultServerSuite) TestSetDefaultServer_ExpectOk() {
	// Given
	serverName := "my-server1"

	s.renderer.EXPECT().RenderServers(gomock.Any()).Times(1)

	// WHEN
	err := s.manager.SetDefaultServer(serverName)

	// THEN
	s.Require().NoError(err)
	kaiConf, err := s.kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(serverName)
	s.Require().NoError(err)
	s.Require().Equal(true, updatedSrv.IsDefault)
}
