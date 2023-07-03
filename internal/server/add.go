package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/spf13/viper"
)

var (
	ErrInvalidKaiServer = errors.New("invalid server")
)

func (c *KaiConfigurator) AddNewServer(server Server) error {
	err := c.validateServer(server)
	if err != nil {
		return fmt.Errorf("validate server: %w", err)
	}

	err = c.addServerToConfiguration(server)

	switch {
	case errors.Is(err, os.ErrNotExist):
		c.logger.Info("Configuration not found, creating a new one.")

		err = c.createInitialConfiguration(server)
		if err != nil {
			return fmt.Errorf("generate initial configuration: %w", err)
		}

	case err != nil:
		return fmt.Errorf("add server to user configuration: %w", err)
	}

	return nil
}

func (c *KaiConfigurator) validateServer(server Server) error {
	if err := server.Validate(); err != nil {
		return err
	}

	if err := c.validateURL(server.URL); err != nil {
		return err
	}

	return nil
}

func (c *KaiConfigurator) validateURL(serverURL string) error {
	// TODO: this enpoint needs to be defined
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

func (c *KaiConfigurator) addServerToConfiguration(server Server) error {
	userConfig, err := c.getConfiguration()
	if err != nil {
		return fmt.Errorf("get user configuration: %w", err)
	}

	if err := userConfig.AddServer(server); err != nil {
		return err
	}

	err = c.writeConfiguration(userConfig)
	if err != nil {
		return fmt.Errorf("update user configuration: %w", err)
	}

	return nil
}

func (c *KaiConfigurator) createInitialConfiguration(server Server) error {
	err := c.createConfigurationDir()
	if err != nil {
		return fmt.Errorf("create configuration directory: %w", err)
	}

	initialConfiguration := &KaiConfiguration{
		DefaultServer: server.Name,
		Servers:       []Server{server},
	}

	err = c.writeConfiguration(initialConfiguration)
	if err != nil {
		return fmt.Errorf("wite user configuration: %w", err)
	}

	return nil
}

func (c *KaiConfigurator) createConfigurationDir() error {
	var (
		configDirPath  = path.Dir(viper.GetString(config.KaiPathKey))
		dirPermissions = 0750
	)

	err := os.Mkdir(configDirPath, os.FileMode(dirPermissions))
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	return nil
}
