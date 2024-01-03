package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/konstellation-io/kli/authserver"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

const (
	_grantTypeRefreshToken       = "refresh_token"
	_refreshTokenRequestTemplate = "%s/realms/%s/protocol/openid-connect/token" //nolint:gosec // False positive
	_logoutRequestTemplate       = "%s/realms/%s/protocol/openid-connect/logout"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials or not found")
)

type AuthenticationService struct {
	logger        logging.Interface
	authServer    authserver.Authenticator
	configService *configuration.KaiConfigService
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
}

func NewAuthentication(logger logging.Interface, authServer authserver.Authenticator) *AuthenticationService {
	return &AuthenticationService{
		logger:        logger,
		authServer:    authServer,
		configService: configuration.NewKaiConfigService(logger),
	}
}

func (a *AuthenticationService) GetToken(serveName string) (*configuration.Token, error) {
	kaiConfig, err := a.configService.GetConfiguration()
	if err != nil {
		return nil, err
	}

	server, err := kaiConfig.GetServer(serveName)
	if err != nil {
		return nil, err
	}

	if !a.areCredentialsValid(server) {
		return nil, ErrInvalidCredentials
	}

	// If the token is valid, return it
	if server.Token != nil && server.Token.IsValid() {
		return server.Token, nil
	}

	// Login to the server
	if server.Token == nil || server.Token.RefreshToken == "" {
		return nil, ErrInvalidToken
	}

	token, err := a.refreshTokenRequest(server)
	if err != nil {
		return nil, err
	}

	if token.AccessToken == "" {
		return nil, ErrInvalidToken
	}

	server.Token = &configuration.Token{
		Date:             time.Now().UTC(),
		AccessToken:      token.AccessToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		RefreshToken:     token.RefreshToken,
		TokenType:        token.TokenType,
	}

	err = kaiConfig.UpdateServer(server)
	if err != nil {
		return nil, err
	}

	err = a.configService.WriteConfiguration(kaiConfig)
	if err != nil {
		return nil, err
	}

	return server.Token, nil
}

func (a *AuthenticationService) Login(serverName, realm, clientID string) (*configuration.Token, error) {
	kaiConfig, err := a.configService.GetConfiguration()
	if err != nil {
		return nil, err
	}

	server, err := kaiConfig.GetServer(serverName)
	if err != nil {
		return nil, err
	}

	// Add credentials to the server
	server.Realm = realm
	server.ClientID = clientID

	// If the credentials are empty, return an error
	if !a.areCredentialsValid(server) {
		return nil, ErrInvalidCredentials
	}

	tokenResponse, err := a.loginRequest(server)
	if err != nil {
		return nil, err
	}

	server.Token = &configuration.Token{
		Date:             time.Now().UTC(),
		AccessToken:      tokenResponse.AccessToken,
		ExpiresIn:        tokenResponse.ExpiresIn,
		RefreshExpiresIn: tokenResponse.RefreshExpiresIn,
		RefreshToken:     tokenResponse.RefreshToken,
		TokenType:        tokenResponse.TokenType,
	}

	err = kaiConfig.UpdateServer(server)
	if err != nil {
		return nil, err
	}

	err = a.configService.WriteConfiguration(kaiConfig)
	if err != nil {
		return nil, err
	}

	return server.Token, nil
}

func (a *AuthenticationService) LoginCLI(serverName, realm, clientID, clientSecret, username,
	password string) (*configuration.Token, error) {
	kaiConfig, err := a.configService.GetConfiguration()
	if err != nil {
		return nil, err
	}

	server, err := kaiConfig.GetServer(serverName)
	if err != nil {
		return nil, err
	}

	// Add credentials to the server
	server.Realm = realm
	server.ClientID = clientID
	server.ClientSecret = clientSecret
	server.Username = username
	server.Password = password

	// If the credentials are empty, return an error
	if !a.areCredentialsValid(server) {
		return nil, ErrInvalidCredentials
	}

	tokenResponse, err := a.loginRequest(server)
	if err != nil {
		return nil, err
	}

	server.Token = &configuration.Token{
		Date:             time.Now().UTC(),
		AccessToken:      tokenResponse.AccessToken,
		ExpiresIn:        tokenResponse.ExpiresIn,
		RefreshExpiresIn: tokenResponse.RefreshExpiresIn,
		RefreshToken:     tokenResponse.RefreshToken,
		TokenType:        tokenResponse.TokenType,
	}

	err = kaiConfig.UpdateServer(server)
	if err != nil {
		return nil, err
	}

	err = a.configService.WriteConfiguration(kaiConfig)
	if err != nil {
		return nil, err
	}

	return server.Token, nil
}

func (a *AuthenticationService) Logout(serverName string) error {
	kaiConfig, err := a.configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := kaiConfig.GetServer(serverName)
	if err != nil {
		return err
	}

	if a.areCredentialsValid(server) {
		err = a.logoutRequest(server)
		if err != nil {
			return err
		}
	}

	server.Realm = ""
	server.ClientID = ""
	server.Token = nil

	err = kaiConfig.UpdateServer(server)
	if err != nil {
		return err
	}

	return a.configService.WriteConfiguration(kaiConfig)
}

func (a *AuthenticationService) loginRequest(server *configuration.Server) (*TokenResponse, error) {
	a.logger.Info("Logging in...")

	authResponse, err := a.authServer.Login(
		&authserver.KeycloakConfig{
			KeycloakURL:  server.AuthEndpoint,
			Realm:        server.Realm,
			ClientID:     server.ClientID,
			Username:     server.Username,
			Password:     server.Password,
			ClientSecret: server.ClientSecret,
		},
	)
	if err != nil {
		return nil, err
	}

	a.logger.Debug(fmt.Sprintf("Login successful. The token is %s", authResponse.AccessToken))

	return &TokenResponse{
		AccessToken:      authResponse.AccessToken,
		ExpiresIn:        authResponse.ExpiresIn,
		RefreshExpiresIn: authResponse.RefreshExpiresIn,
		RefreshToken:     authResponse.RefreshToken,
		TokenType:        authResponse.TokenType,
	}, nil
}

func (a *AuthenticationService) logoutRequest(server *configuration.Server) error {
	a.logger.Info("Logging out...")

	u, err := url.Parse(fmt.Sprintf(_logoutRequestTemplate, server.AuthEndpoint, server.Realm))
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("client_id", server.ClientID)
	data.Add("refresh_token", server.Token.RefreshToken)

	// Make the HTTP POST request
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", server.Token.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error requesting token: %s", resp.Status)
	}

	return nil
}

func (a *AuthenticationService) refreshTokenRequest(server *configuration.Server) (*TokenResponse, error) {
	a.logger.Info("Refreshing token...")

	u, err := url.Parse(fmt.Sprintf(_refreshTokenRequestTemplate, server.AuthEndpoint, server.Realm))
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("client_id", server.ClientID)
	data.Set("grant_type", _grantTypeRefreshToken)
	data.Add("refresh_token", server.Token.RefreshToken)

	// Make the HTTP POST request
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refreshing token: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)

	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (a *AuthenticationService) areCredentialsValid(server *configuration.Server) bool {
	return server.AuthEndpoint != "" && server.Realm != "" && server.ClientID != ""
}
