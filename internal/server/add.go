package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/konstellation-io/kli/internal/configuration"
)

var (
	_validServerName     = regexp.MustCompile(`^[A-Za-z0-9\-_]+$`)
	ErrInvalidKaiServer  = errors.New("invalid server")
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

	kaiConfig, err := c.configHandler.GetConfiguration()
	if err != nil {
		return err
	}

	c.renderer.RenderServers(kaiConfig.Servers)

	return nil
}

func (c *Handler) validateServer(server Server) error {
	if err := c.validateServerName(server); err != nil {
		return err
	}

	return c.validateURL(server.URL)
}

func (c *Handler) validateServerName(server Server) error {
	if !_validServerName.MatchString(server.Name) {
		return ErrInvalidServerName
	}

	return nil
}

func (c *Handler) validateURL(serverURL string) error {
	infoURL, err := url.JoinPath(serverURL, "info")
	if err != nil {
		return fmt.Errorf("format URL: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, infoURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("initialize request to %s: %w", infoURL, err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error querying %s: %w", infoURL, err)
	}

	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	var serverInfo RemoteServerInfo

	err = json.Unmarshal(responseData, &serverInfo)
	if err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	if !serverInfo.IsKaiServer {
		return ErrInvalidKaiServer
	}

	return nil
}

func (c *Handler) addServerToConfiguration(server Server, isDefault bool) error {
	userConfig, err := c.configHandler.GetConfiguration()
	if err != nil {
		return fmt.Errorf("get user configuration: %w", err)
	}

	if err := userConfig.AddServer(configuration.Server{
		Name:      server.Name,
		URL:       server.URL,
		IsDefault: isDefault,
	}); err != nil {
		return err
	}

	err = c.configHandler.WriteConfiguration(userConfig)
	if err != nil {
		return fmt.Errorf("update user configuration: %w", err)
	}

	return nil
}
