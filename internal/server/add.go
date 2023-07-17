package server

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/konstellation-io/kli/internal/configuration"
)

var (
	_validServerName     = regexp.MustCompile(`^[A-Za-z0-9\-_]+$`)
	ErrInvalidServerName = errors.New("invalid server name, only alphanumeric and hyphens are supported")
)

func (c *Handler) AddNewServer(server Server, isDefault bool) error {
	err := c.validateServer(server)
	if err != nil {
		return fmt.Errorf("validate server: %w", err)
	}

	err = c.addServerToConfiguration(server, isDefault)

	if err != nil {
		return fmt.Errorf("add server to user configuration: %w", err)
	}

	kaiConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return err
	}

	c.renderer.RenderServers(kaiConfig.Servers)

	return nil
}

func (c *Handler) validateServer(server Server) error {
	return c.validateServerName(server)
}

func (c *Handler) validateServerName(server Server) error {
	if !_validServerName.MatchString(server.Name) {
		return ErrInvalidServerName
	}

	return nil
}

func (c *Handler) addServerToConfiguration(server Server, isDefault bool) error {
	userConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return fmt.Errorf("get user configuration: %w", err)
	}

	if err := userConfig.AddServer(&configuration.Server{
		Name:      server.Name,
		URL:       server.URL,
		IsDefault: isDefault,
	}); err != nil {
		return err
	}

	err = c.configService.WriteConfiguration(userConfig)
	if err != nil {
		return fmt.Errorf("update user configuration: %w", err)
	}

	return nil
}
