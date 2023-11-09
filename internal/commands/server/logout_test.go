package server_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ServerLogoutSuite struct {
	suite.Suite

	renderer *mocks.MockRenderer
	manager  *server.Handler
	logger   logging.Interface
	tmpDir   string
}

func TestServerLogoutSuite(t *testing.T) {
	suite.Run(t, new(ServerLogoutSuite))
}

func (s *ServerLogoutSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.logger = logger
	s.renderer = renderer
	s.manager = server.NewHandler(logger, renderer)

	tmpDir, err := os.MkdirTemp("", "TestServerLogout_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.tmpDir = tmpDir
}

func (s *ServerLogoutSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ServerLogoutSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ServerLogoutSuite) TestLogoutServer_ExpectOk() {
	// Given
	err := os.MkdirAll(path.Dir(viper.GetString(config.KaiConfigPath)), 0750)
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConfService := configuration.NewKaiConfigService(s.logger)

	kaiConf, err := kaiConfService.GetConfiguration()
	s.Require().NoError(err)

	srv := &configuration.Server{
		Name:      "my-server",
		Host:      "https://kai-dev.konstellation.io",
		AuthURL:   "https://auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		IsDefault: true,
		Token:     &configuration.Token{},
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout",
		srv.AuthURL, srv.Realm),
		httpmock.NewStringResponder(http.StatusNoContent, `{}`))

	// WHEN
	err = s.manager.Logout(srv.Name)

	// THEN
	s.Require().NoError(err)
	kaiConf, err = kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(srv.Name)
	s.Require().NoError(err)
	s.Require().Empty(updatedSrv.AuthURL)
	s.Require().Empty(updatedSrv.Realm)
	s.Require().Empty(updatedSrv.ClientID)
	s.Require().Nil(updatedSrv.Token)
}

func (s *ServerLogoutSuite) TestLogoutServer_NotLoggedInServer_ExpectOk() {
	// Given
	err := os.MkdirAll(path.Dir(viper.GetString(config.KaiConfigPath)), 0750)
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConfService := configuration.NewKaiConfigService(s.logger)

	kaiConf, err := kaiConfService.GetConfiguration()
	s.Require().NoError(err)

	srv := &configuration.Server{
		Name:      "my-server",
		Host:      "https://kai-dev.konstellation.io",
		IsDefault: true,
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout",
		srv.AuthURL, srv.Realm),
		httpmock.NewStringResponder(http.StatusUnauthorized, `{}`))

	// WHEN
	err = s.manager.Logout(srv.Name)

	// THEN
	s.Require().NoError(err)
	kaiConf, err = kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(srv.Name)
	s.Require().NoError(err)
	s.Require().Empty(updatedSrv.AuthURL)
	s.Require().Empty(updatedSrv.Realm)
	s.Require().Empty(updatedSrv.ClientID)
	s.Require().Nil(updatedSrv.Token)
}
