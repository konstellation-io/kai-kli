package server

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Handler) GetToken(serverName string) (*configuration.Token, error) {
	return c.authentication.GetToken(serverName)
}
