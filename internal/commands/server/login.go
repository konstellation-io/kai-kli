package server

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Handler) Login(serverName, authURL, realm, clientID, username, password string) (*configuration.Token, error) {
	return c.authentication.Login(serverName, authURL, realm, clientID, username, password)
}
