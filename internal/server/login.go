package server

import (
	"errors"

	"github.com/konstellation-io/kli/internal/configuration"
)

var (
	ErrCouldNotLogin = errors.New("could not login")
)

func (c *Handler) Login(serverName, realm, clientID, username, password string) (*configuration.Token, error) {
	return c.authentication.Login(serverName, realm, clientID, username, password)
}
