package authserver

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/konstellation-io/kli/internal/logging"
)

type CLIAuthenticator struct {
	logger logging.Interface
	client *gocloak.GoCloak
}

func NewCLIAuthenticator(logger logging.Interface, client *gocloak.GoCloak) *CLIAuthenticator {
	return &CLIAuthenticator{
		logger: logger,
		client: client,
	}
}

func (a *CLIAuthenticator) Login(config *KeycloakConfig) (*AuthResponse, error) {
	token, err := a.client.Login(
		context.Background(),
		config.ClientID,
		config.ClientSecret,
		config.Realm,
		config.Username,
		config.Password,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:      token.AccessToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		RefreshToken:     token.RefreshToken,
		TokenType:        token.TokenType,
	}, nil
}
