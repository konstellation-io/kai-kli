package auth_test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	auth "github.com/konstellation-io/kli/internal/authentication"
	"github.com/konstellation-io/kli/internal/configuration"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/mocks"
)

type AuthenticationSuite struct {
	suite.Suite

	authentication *auth.Authentication
	logger         logging.Interface
	tmpDir         string
}

func TestAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationSuite))
}

func (s *AuthenticationSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	s.logger = logger
	s.authentication = auth.NewAuthentication(logger)

	tmpDir, err := os.MkdirTemp("", "TestAuthentication_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.tmpDir = tmpDir
}

func (s *AuthenticationSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *AuthenticationSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *AuthenticationSuite) TestLogin_ExpectToken() {
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
		URL:       "auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
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
	httpmock.RegisterMatcherResponder("POST", fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token",
		srv.URL, srv.Realm),
		httpmock.HeaderExists("Content-Type").
			And(httpmock.HeaderContains("Content-Type", "application/x-www-form-urlencoded")),
		httpmock.NewStringResponder(200, `{
					"access_token": "access token",
					"expires_in": 300,
					"refresh_expires_in": 5000,
					"refresh_token": "refresh token",
					"token_type": "Bearer"}`))

	// WHEN
	token, err := s.authentication.Login(srv.Name, srv.Realm, srv.ClientID, srv.Username, srv.Password)

	// THEN
	s.Require().NoError(err)
	s.Require().NotNil(token)
	kaiConf, err = kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(srv.Name)
	s.Require().NoError(err)
	s.Require().Equal("access token", updatedSrv.Token.AccessToken)
}

func (s *AuthenticationSuite) TestLogin_ExpectError() {
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
		URL:       "auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
		IsDefault: true,
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token",
		srv.URL, srv.Realm),
		httpmock.NewErrorResponder(errors.New("error getting token")))

	// WHEN
	token, err := s.authentication.Login(srv.Name, srv.Realm, srv.ClientID, srv.Username, srv.Password)

	// THEN
	s.Require().Error(err)
	s.Require().Nil(token)
}

func (s *AuthenticationSuite) TestGetToken_GetNewToken_ExpectToken() {
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
		URL:       "auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
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
	httpmock.RegisterMatcherResponder("POST", fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token",
		srv.URL, srv.Realm),
		httpmock.HeaderExists("Content-Type").
			And(httpmock.HeaderContains("Content-Type", "application/x-www-form-urlencoded")),
		httpmock.NewStringResponder(200, `{
					"access_token": "new access token",
					"expires_in": 300,
					"refresh_expires_in": 5000,
					"refresh_token": "refresh token",
					"token_type": "Bearer"}`))

	// WHEN
	token, err := s.authentication.GetToken(srv.Name)

	// THEN
	s.Require().NoError(err)
	s.Require().NotNil(token)

	kaiConf, err = kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(srv.Name)
	s.Require().NoError(err)
	s.Require().Equal("new access token", updatedSrv.Token.AccessToken)
}

func (s *AuthenticationSuite) TestGetToken_GetStoredToken_ExpectToken() {
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
		URL:       "auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
		IsDefault: true,
		Token: &configuration.Token{
			Date:             time.Now().UTC(),
			AccessToken:      "stored token",
			ExpiresIn:        3600,
			RefreshExpiresIn: 3600,
			RefreshToken:     "refresh token",
			TokenType:        "Bearer",
		},
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	// WHEN
	token, err := s.authentication.GetToken(srv.Name)

	// THEN
	s.Require().NoError(err)
	s.Require().NotNil(token)
	kaiConf, err = kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(srv.Name)
	s.Require().NoError(err)
	s.Require().Equal("stored token", updatedSrv.Token.AccessToken)
}

func (s *AuthenticationSuite) TestGetToken_GetRenewedToken_ExpectToken() {
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
		URL:       "auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
		IsDefault: true,
		Token: &configuration.Token{
			Date:             time.Now().AddDate(0, 0, -1).UTC(),
			AccessToken:      "stored token",
			ExpiresIn:        3600,
			RefreshExpiresIn: 3600,
			RefreshToken:     "refresh token",
			TokenType:        "Bearer",
		},
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterMatcherResponder("POST", fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token",
		srv.URL, srv.Realm),
		httpmock.HeaderExists("Content-Type").
			And(httpmock.HeaderContains("Content-Type", "application/x-www-form-urlencoded")),
		httpmock.NewStringResponder(200, `{
					"access_token": "renewed access token",
					"expires_in": 300,
					"refresh_expires_in": 5000,
					"refresh_token": "refresh token",
					"token_type": "Bearer"}`))

	// WHEN
	token, err := s.authentication.GetToken(srv.Name)

	// THEN
	s.Require().NoError(err)
	s.Require().NotNil(token)
	kaiConf, err = kaiConfService.GetConfiguration()
	s.Require().NoError(err)
	updatedSrv, err := kaiConf.GetServer(srv.Name)
	s.Require().NoError(err)
	s.Require().Equal("renewed access token", updatedSrv.Token.AccessToken)
}

func (s *AuthenticationSuite) TestGetToken_ExpectError() {
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
		URL:       "auth.kai-dev.konstellation.io",
		Realm:     "konstellation",
		ClientID:  "admin-cli",
		Username:  "david",
		Password:  "password",
		IsDefault: true,
	}

	err = kaiConf.AddServer(srv)
	s.Require().NoError(err)

	err = kaiConfService.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token",
		srv.URL, srv.Realm),
		httpmock.NewErrorResponder(errors.New("error getting token")))

	// WHEN
	token, err := s.authentication.GetToken(srv.Name)

	// THEN
	s.Require().Error(err)
	s.Require().Nil(token)
}
