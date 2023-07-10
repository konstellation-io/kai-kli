package configuration_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ConfigurationHandlerTest struct {
	suite.Suite

	configHandler configuration.KaiConfigHandler
	tmpDir        string
}

func TestConfigurationHandlerSuite(t *testing.T) {
	suite.Run(t, new(ConfigurationHandlerTest))
}

func (ch *ConfigurationHandlerTest) SetupSuite() {
	ctrl := gomock.NewController(ch.T())
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	ch.configHandler = *configuration.NewKaiConfigHandler(logger)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	ch.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	ch.tmpDir = tmpDir
}

func (ch *ConfigurationHandlerTest) TearDownSuite(_, _ string) {
	err := os.RemoveAll(ch.tmpDir)
	ch.Require().NoError(err)
}

func (ch *ConfigurationHandlerTest) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			ch.T().Fatalf("error cleaning tmp path: %ch", err)
		}
	}
}

func (ch *ConfigurationHandlerTest) TestReadConfig_ReadExistingConfiguration_ExpectOK() {
	// GIVEN
	defaultConfig, err := createDefaultConfiguration()
	ch.Require().NoError(err)

	// WHEN
	currentConfig, err := ch.configHandler.GetConfiguration()

	// THEN
	ch.Require().NoError(err)
	ch.NotNil(currentConfig)
	ch.Equal(defaultConfig, currentConfig)
}

func (ch *ConfigurationHandlerTest) TestReadConfig_ReadNonExistingConfiguration_ExpectError() {
	// GIVEN no config file
	// WHEN
	currentConfig, err := ch.configHandler.GetConfiguration()

	// THEN
	ch.Require().Error(err)
	ch.Nil(currentConfig)
}

func (ch *ConfigurationHandlerTest) TestWriteConfig_WriteValidConfiguration_ExpectOK() {
	// GIVEN
	defaultConfig := getDefaultConfiguration()
	err := os.MkdirAll(path.Dir(viper.GetString(config.KaiConfigPath)), 0750)
	ch.Require().NoError(err)
	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	ch.Require().NoError(err)

	// WHEN
	err = ch.configHandler.WriteConfiguration(&defaultConfig)

	// THEN
	ch.Require().NoError(err)
	currentConfig, err := ch.configHandler.GetConfiguration()
	ch.Require().NoError(err)
	ch.NotNil(currentConfig)
	ch.Equal(defaultConfig, *currentConfig)
}

func (ch *ConfigurationHandlerTest) TestWriteConfig_WriteInvalidConfiguration_ExpectError() {
	// GIVEN a nil configuration
	// WHEN
	err := ch.configHandler.WriteConfiguration(nil)

	// THEN
	ch.Require().Error(err)
}

func getDefaultConfiguration() configuration.KaiConfiguration {
	return configuration.KaiConfiguration{
		Servers: []configuration.Server{
			{
				Name:      "existing-server",
				URL:       "existing-server.com",
				IsDefault: true,
			},
		},
	}
}

func createDefaultConfiguration() (*configuration.KaiConfiguration, error) {
	defaultConfig := getDefaultConfiguration()

	if err := createConfigWithServer(defaultConfig); err != nil {
		return nil, err
	}

	return &defaultConfig, nil
}

func createConfigWithServer(cfg configuration.KaiConfiguration) error {
	if err := os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), 0750); err != nil {
		return err
	}

	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(viper.GetString(config.KaiConfigPath), configBytes, 0600)
}
