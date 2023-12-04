package server

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/konstellation-io/kli/internal/services/configuration"
)

const (
	_apiSubdomain     = "api"
	_authSubdomain    = "auth"
	_storageSubdomain = "storage-console"
)

var (
	_validServerName = regexp.MustCompile(`^[A-Za-z0-9\-_]+$`)
	_validHostName   = regexp.MustCompile(
		`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)

	ErrInvalidServerName = errors.New("invalid server name, only alphanumeric and hyphens are supported")
	ErrInvalidServerHost = errors.New("invalid server host, make sure you are NOT including protocol scheme")
)

func (c *Handler) AddNewServer(server *Server, isDefault bool) error {
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

func (c *Handler) validateServer(server *Server) error {
	if err := c.validateServerName(server); err != nil {
		return err
	}

	return c.validateServerHost(server)
}

func (c *Handler) validateServerName(server *Server) error {
	if !_validServerName.MatchString(server.Name) {
		return ErrInvalidServerName
	}

	return nil
}

func (c *Handler) validateServerHost(server *Server) error {
	if !_validHostName.MatchString(server.Host) {
		return ErrInvalidServerHost
	}

	return nil
}

func (c *Handler) addServerToConfiguration(server *Server, isDefault bool) error {
	userConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return fmt.Errorf("get user configuration: %w", err)
	}

	if err := userConfig.AddServer(&configuration.Server{
		Name:            server.Name,
		Host:            server.Host,
		Protocol:        server.Protocol,
		APIEndpoint:     c.getAPIEndpoint(server.Protocol, server.Host),
		AuthEndpoint:    c.getAuthEndpoint(server.Protocol, server.Host),
		StorageEndpoint: c.getStorageEndpoint(server.Protocol, server.Host),
		IsDefault:       isDefault,
	}); err != nil {
		return err
	}

	err = c.configService.WriteConfiguration(userConfig)
	if err != nil {
		return fmt.Errorf("update user configuration: %w", err)
	}

	return nil
}

func (c *Handler) getAPIEndpoint(protocol, host string) string {
	return fmt.Sprintf("%s://%s.%s", protocol, _apiSubdomain, host)
}

func (c *Handler) getAuthEndpoint(protocol, host string) string {
	return fmt.Sprintf("%s://%s.%s", protocol, _authSubdomain, host)
}

func (c *Handler) getStorageEndpoint(protocol, host string) string {
	return fmt.Sprintf("%s://%s.%s", protocol, _storageSubdomain, host)
}
