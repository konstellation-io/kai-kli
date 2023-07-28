package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

const (
	DebugKey             = "debug"
	RequestTimeoutKey    = "request_timeout"
	BuildVersionKey      = "build_version"
	KaiConfigPath        = "kai_path"
	KaiConfigFile        = "kai_config_filename"
	KaiProductConfigFile = "kai_product_config_filename"

	_defaultKaiConfigFile        = "kai.conf"
	_defaultKaiProductConfigFile = "krt.yaml"
	_defaultKaiFolder            = ".kai"
	_defaultRequestTimeout       = 2 * time.Minute
)

func InitConfigWithBuildVersion(buildVersion string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get user home path: %w", err)
	}

	viper.SetDefault(KaiConfigFile, _defaultKaiConfigFile)
	viper.SetDefault(KaiConfigPath, path.Join(userHomeDir, _defaultKaiFolder, _defaultKaiConfigFile))
	viper.SetDefault(KaiProductConfigFile, _defaultKaiProductConfigFile)
	viper.Set(BuildVersionKey, buildVersion)
	viper.Set(RequestTimeoutKey, _defaultRequestTimeout)
	viper.Set(DebugKey, false)

	return nil
}

func GetProductConfigFilePath(product string, productPath ...string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current directory: %w", err)
	}

	productConfigPath := currentDir

	if len(productPath) != 0 {
		productConfigPath = path.Join(productPath...)
	}

	return path.Join(productConfigPath, _defaultKaiFolder, product, viper.GetString(KaiProductConfigFile)), nil
}
