package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidKaiServer     = errors.New("invalid kai server")
	ErrDuplicatedServerName = errors.New("duplicated server name")
	ErrDuplicatedServerURL  = errors.New("duplicated server URL")
)

type KaiConfiguration struct {
	Servers []Server
}

func (kc *KaiConfiguration) AddServer(server Server) error {
	if err := kc.checkServerDuplication(server); err != nil {
		return err
	}

	kc.Servers = append(kc.Servers, server)

	return nil
}

func (kc *KaiConfiguration) checkServerDuplication(server Server) error {
	for _, s := range kc.Servers {
		if s.Name == server.Name {
			return ErrDuplicatedServerName
		}

		if s.URL == server.URL {
			return ErrDuplicatedServerURL
		}
	}

	return nil
}

type Manager struct {
	logger logging.Interface
}

func NewManager(logger logging.Interface) *Manager {
	return &Manager{
		logger: logger,
	}
}

func (sm *Manager) AddNewServer(server Server) error {
	err := sm.validateServer(server)
	if err != nil {
		return fmt.Errorf("validate server: %w", err)
	}

	err = sm.addServerToConfiguration(server)

	switch {
	case errors.Is(err, os.ErrNotExist):
		sm.logger.Info("Configuration not found, creating a new one.")

		err = sm.createInitialConfiguration(server)
		if err != nil {
			return fmt.Errorf("generate initial configuration: %w", err)
		}

	case err != nil:
		return fmt.Errorf("add server to user configuration: %w", err)
	}

	return nil
}

func (sm *Manager) validateServer(server Server) error {
	if err := server.Validate(); err != nil {
		return err
	}

	if err := sm.validateURL(server.URL); err != nil {
		return err
	}

	return nil
}

func (sm *Manager) validateURL(url string) error {
	// TODO: this enpoint needs to be defined
	response, err := http.Get(url + "/info")
	if err != nil {
		return err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var serverInfo RemoteServerInfo

	err = json.Unmarshal(responseData, &serverInfo)
	if err != nil {
		return err
	}

	if !serverInfo.IsKaiServer {
		return errors.New("invalid server")
	}

	return nil
}

func (sm *Manager) addServerToConfiguration(server Server) error {
	userConfig, err := sm.getKaiConfiguration()
	if err != nil {
		return fmt.Errorf("get user configuration: %w", err)
	}

	if err := userConfig.AddServer(server); err != nil {
		return err
	}

	err = sm.writeConfiguration(userConfig)
	if err != nil {
		return fmt.Errorf("update user configuration: %w", err)
	}

	return nil
}

func (sm *Manager) getKaiConfiguration() (*KaiConfiguration, error) {
	configBytes, err := os.ReadFile(viper.GetString(config.KaiPathKey))
	if err != nil {
		return nil, err
	}

	var kaiConfiguration KaiConfiguration

	err = yaml.Unmarshal(configBytes, &kaiConfiguration)
	if err != nil {
		return nil, err
	}

	return &kaiConfiguration, nil
}

func (sm *Manager) writeConfiguration(newConfig *KaiConfiguration) error {
	updatedConfig, err := yaml.Marshal(newConfig)
	if err != nil {
		return err
	}

	filePermissions := 0600

	err = os.WriteFile(viper.GetString(config.KaiPathKey), updatedConfig, os.FileMode(filePermissions))
	if err != nil {
		return err
	}

	return nil
}

func (sm *Manager) createInitialConfiguration(server Server) error {
	err := sm.createConfigurationDir()
	if err != nil {
		return fmt.Errorf("create configuration directory: %w", err)
	}

	server.Default = true

	initialConfiguration := &KaiConfiguration{
		Servers: []Server{server},
	}

	err = sm.writeConfiguration(initialConfiguration)
	if err != nil {
		return fmt.Errorf("wite user configuration: %w", err)
	}

	return nil
}

func (sm *Manager) createConfigurationDir() error {
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
