package configuration

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/logging"
)

type KaiConfigHandler struct {
	logger logging.Interface
}

func NewKaiConfigHandler(logger logging.Interface) *KaiConfigHandler {
	return &KaiConfigHandler{
		logger: logger,
	}
}

func (c *KaiConfigHandler) GetConfiguration() (*KaiConfiguration, error) {
	configBytes, err := os.ReadFile(viper.GetString(config.KaiConfigPath))
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

func (c *KaiConfigHandler) WriteConfiguration(newConfig *KaiConfiguration) error {
	updatedConfig, err := yaml.Marshal(newConfig)
	if err != nil {
		return err
	}

	filePermissions := 0600

	err = os.WriteFile(viper.GetString(config.KaiConfigPath), updatedConfig, os.FileMode(filePermissions))
	if err != nil {
		return err
	}

	return nil
}
