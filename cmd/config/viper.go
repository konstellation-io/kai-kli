package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

const (
	DebugKey               = "debug"
	RequestTimeoutKey      = "request_timeout"
	BuildVersionKey        = "build_version"
	KaiConfigPath          = "kai_path"
	KaiConfigFile          = "kai_config_filename"
	KaiProductConfigFolder = "kai_product_config_filename"

	_defaultKaiConfigFile  = "kai.conf"
	_defaultKaiFolder      = ".kai"
	_defaultRequestTimeout = 5 * time.Minute
)

func InitConfigWithBuildVersion(buildVersion string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get user home path: %w", err)
	}

	viper.SetDefault(KaiConfigFile, _defaultKaiConfigFile)
	viper.SetDefault(KaiConfigPath, path.Join(userHomeDir, _defaultKaiFolder, _defaultKaiConfigFile))
	viper.SetDefault(KaiProductConfigFolder, _defaultKaiFolder)
	viper.Set(BuildVersionKey, buildVersion)
	viper.Set(RequestTimeoutKey, _defaultRequestTimeout)
	viper.Set(DebugKey, false)

	return nil
}
