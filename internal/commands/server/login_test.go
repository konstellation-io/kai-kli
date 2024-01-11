//go:build unit

package server_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ServerLoginSuite struct {
	suite.Suite

	renderer      *mocks.MockRenderer
	manager       *server.Handler
	logger        logging.Interface
	tmpDir        string
	authenticator *mocks.MockAuthenticator
}

func TestServerLoginSuite(t *testing.T) {
	suite.Run(t, new(ServerLoginSuite))
}

func (s *ServerLoginSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.logger = logger
	s.renderer = renderer
	s.authenticator = mocks.NewMockAuthenticator(ctrl)
	s.manager = server.NewTestHandler(logger, renderer, s.authenticator)

	tmpDir, err := os.MkdirTemp("", "TestServerLogin_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.tmpDir = tmpDir
}

func (s *ServerLoginSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ServerLoginSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ServerLoginSuite) TestLoginServer_ExpectToken() {
	// Given
	err := os.MkdirAll(path.Dir(viper.GetString(config.KaiConfigPath)), 0750)
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConfService := configuration.NewKaiConfigService(s.logger)

	kaiConf, err := kaiConfService.GetConfiguration()
	s.Require().NoError(err)

	srv := &configuration.Server{
		Name:         "my-server",
		Host:         "https://kai-dev.konstellation.io",
		AuthEndpoint: "https://auth.kai-dev.konstellation.io",
		Realm:        "konstellation",
		ClientID:     "admin-cli",
		IsDefault:    true,
		Token:        &configuration.Token{},
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	s.authenticator.EXPECT().Login(gomock.Any()).Return(&authserver.AuthResponse{
		AccessToken:      "access token",
		ExpiresIn:        300,
		RefreshExpiresIn: 5000,
		RefreshToken:     "refresh token",
		TokenType:        "Bearer",
	}, nil)

	// WHEN
	token, err := s.manager.Login(srv.Name, srv.Realm, srv.ClientID)

	// THEN
	s.Require().NotNil(token)
	s.Require().NoError(err)
	kaiConf, err = kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(srv.Name)
	s.Require().NoError(err)
	s.Require().Equal("access token", updatedSrv.Token.AccessToken)
}

func (s *ServerLoginSuite) TestLoginServer_ExpectError() {
	// Given
	err := os.MkdirAll(path.Dir(viper.GetString(config.KaiConfigPath)), 0750)
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConfService := configuration.NewKaiConfigService(s.logger)

	kaiConf, err := kaiConfService.GetConfiguration()
	s.Require().NoError(err)

	srv := &configuration.Server{
		Name:         "my-server",
		Host:         "https://kai-dev.konstellation.io",
		AuthEndpoint: "https://auth.kai-dev.konstellation.io",
		Realm:        "konstellation",
		ClientID:     "admin-cli",
		IsDefault:    true,
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	s.authenticator.EXPECT().Login(gomock.Any()).Return(nil, errors.New("error getting token"))

	// WHEN
	token, err := s.manager.Login(srv.Name, srv.Realm, srv.ClientID)

	// THEN
	s.Require().Error(err)
	s.Require().Nil(token)
}
