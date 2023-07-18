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

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

const (
	_grantTypePassword    = "password"
	_loginRequestTemplate = "https://%s/realms/%s/protocol/openid-connect/token"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials or not found")
)

type AuthenticationService struct {
	logger        logging.Interface
	configService *configuration.KaiConfigService
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
}

func NewAuthentication(logger logging.Interface) *AuthenticationService {
	return &AuthenticationService{
		logger:        logger,
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
	token, err := a.Login(server.Name, server.AuthURL, server.Realm, server.ClientID, server.Username, server.Password)
	if err != nil {
		return nil, err
	}

	if token.AccessToken == "" {
		return nil, ErrInvalidToken
	}

	return token, nil
}

func (a *AuthenticationService) Login(serverName, authURL, realm, clientID, username, password string) (*configuration.Token, error) {
	kaiConfig, err := a.configService.GetConfiguration()
	if err != nil {
		return nil, err
	}

	server, err := kaiConfig.GetServer(serverName)
	if err != nil {
		return nil, err
	}

	// Add credentials to the server
	server.AuthURL = authURL
	server.Realm = realm
	server.ClientID = clientID
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

func (a *AuthenticationService) loginRequest(server *configuration.Server) (*TokenResponse, error) {
	u, err := url.Parse(fmt.Sprintf(_loginRequestTemplate, server.AuthURL, server.Realm))
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("username", server.Username)
	data.Add("password", server.Password)
	data.Add("grant_type", _grantTypePassword)
	data.Add("client_id", server.ClientID)

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
		return nil, fmt.Errorf("error requesting token: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)

	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (a *AuthenticationService) areCredentialsValid(server *configuration.Server) bool {
	return server.AuthURL != "" && server.Username != "" &&
		server.Password != "" && server.Realm != "" && server.ClientID != ""
}
