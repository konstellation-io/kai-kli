package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

const (
	DebugKey          = "debug"
	RequestTimeoutKey = "request_timeout"
	BuildVersionKey   = "build_version"
	KaiPathKey        = "kai_path"

	_defaultRequestTimeout = 2 * time.Minute
)

func InitConfigWithBuildVersion(buildVersion string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get user home path: %w", err)
	}

	viper.SetDefault(KaiPathKey, path.Join(userHomeDir, ".kai", "kai.conf"))
	viper.Set(BuildVersionKey, buildVersion)
	viper.Set(RequestTimeoutKey, _defaultRequestTimeout)
	viper.Set(DebugKey, false)

	return nil
}
