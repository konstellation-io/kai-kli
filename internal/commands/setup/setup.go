package setup

import (
	"errors"
	"os"
	"path"

	"github.com/spf13/viper"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/logging"
)

type KaiSetup struct {
	logger logging.Interface
}

func NewHandler(logger logging.Interface) *KaiSetup {
	return &KaiSetup{
		logger: logger,
	}
}

func (c *KaiSetup) CreateConfiguration() error {
	err := c.createConfigurationDir()
	if err != nil {
		return err
	}

	err = c.createConfigurationFile()
	if err != nil {
		return err
	}

	return nil
}

func (c *KaiSetup) CheckConfigurationExists() bool {
	return fileExists(viper.GetString(config.KaiConfigPath))
}

func (c *KaiSetup) createConfigurationDir() error {
	var (
		configDirPath  = path.Dir(viper.GetString(config.KaiConfigPath))
		dirPermissions = 0750
	)

	err := os.MkdirAll(configDirPath, os.FileMode(dirPermissions))
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	return nil
}

func (c *KaiSetup) createConfigurationFile() error {
	if fileExists(viper.GetString(config.KaiConfigPath)) {
		return nil
	}

	c.logger.Info("The KAI client is not setup, initializing...")

	_, e := os.Create(viper.GetString(config.KaiConfigPath))

	return e
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
