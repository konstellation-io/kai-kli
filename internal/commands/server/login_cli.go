package server

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Handler) LoginCLI(serverName, realm, clientID, clientSecret, username, password string) (*configuration.Token, error) {
	return c.authentication.LoginCLI(serverName, realm, clientID, clientSecret, username, password)
}
