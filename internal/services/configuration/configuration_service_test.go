package configuration_test

import (
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ConfigurationServiceTest struct {
	suite.Suite

	configService configuration.KaiConfigService
	tmpDir        string
}

func TestConfigurationServiceSuite(t *testing.T) {
	suite.Run(t, new(ConfigurationServiceTest))
}

func (ch *ConfigurationServiceTest) SetupSuite() {
	ctrl := gomock.NewController(ch.T())
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	ch.configService = *configuration.NewKaiConfigService(logger)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	ch.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	ch.tmpDir = tmpDir
}

func (ch *ConfigurationServiceTest) TearDownSuite(_, _ string) {
	err := os.RemoveAll(ch.tmpDir)
	ch.Require().NoError(err)
}

func (ch *ConfigurationServiceTest) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			ch.T().Fatalf("error cleaning tmp path: %ch", err)
		}
	}
}

func (ch *ConfigurationServiceTest) TestReadConfig_ReadExistingConfiguration_ExpectOK() {
	// GIVEN
	defaultConfig, err := createDefaultConfiguration()
	ch.Require().NoError(err)

	// WHEN
	currentConfig, err := ch.configService.GetConfiguration()

	// THEN
	ch.Require().NoError(err)
	ch.NotNil(currentConfig)
	ch.Equal(defaultConfig, currentConfig)
}

func (ch *ConfigurationServiceTest) TestReadConfig_ReadNonExistingConfiguration_ExpectError() {
	// GIVEN no config file
	// WHEN
	currentConfig, err := ch.configService.GetConfiguration()

	// THEN
	ch.Require().Error(err)
	ch.Nil(currentConfig)
}

func (ch *ConfigurationServiceTest) TestWriteConfig_WriteValidConfiguration_ExpectOK() {
	// GIVEN
	defaultConfig := getDefaultConfiguration()
	err := os.MkdirAll(path.Dir(viper.GetString(config.KaiConfigPath)), 0750)
	ch.Require().NoError(err)
	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	ch.Require().NoError(err)

	// WHEN
	err = ch.configService.WriteConfiguration(&defaultConfig)

	// THEN
	ch.Require().NoError(err)
	currentConfig, err := ch.configService.GetConfiguration()
	ch.Require().NoError(err)
	ch.NotNil(currentConfig)
	ch.Equal(defaultConfig, *currentConfig)
}

func (ch *ConfigurationServiceTest) TestWriteConfig_WriteInvalidConfiguration_ExpectError() {
	// GIVEN a nil configuration
	// WHEN
	err := ch.configService.WriteConfiguration(nil)

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
				Realm:     "existing-realm",
				ClientID:  "existing-client-id",
				Username:  "user",
				Password:  "pass",
				Token: &configuration.Token{
					Date:             time.Now().UTC(),
					AccessToken:      "access-token",
					ExpiresIn:        3600,
					RefreshExpiresIn: 3600,
					RefreshToken:     "refresh-token",
					TokenType:        "bearer",
				},
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
