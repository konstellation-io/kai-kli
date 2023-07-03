package server

import (
	"os"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type KaiConfigurator struct {
	logger logging.Interface
}

func NewKaiConfigurator(logger logging.Interface) *KaiConfigurator {
	return &KaiConfigurator{
		logger: logger,
	}
}

func (c *KaiConfigurator) getConfiguration() (*KaiConfiguration, error) {
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

func (c *KaiConfigurator) writeConfiguration(newConfig *KaiConfiguration) error {
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
