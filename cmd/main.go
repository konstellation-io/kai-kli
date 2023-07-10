package main

import (
	"fmt"
	"log"

	"github.com/konstellation-io/kli/api/kai/config"
	config2 "github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/cmd/root"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/pkg/iostreams"
)

// version is dynamically set by the toolchain or overridden at build time.
var version = "DEV"

// date is dynamically set at build time with format YYYY-MM-DD.
var date = "" //nolint:gochecknoglobals

func main() {
	logger := logging.NewDefaultLogger()
	io := iostreams.System()

	cfg, err := config.NewConfig(version)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = config2.InitConfigWithBuildVersion(version)
	if err != nil {
		log.Fatalf("initialize configuration: %s", err)
	}

	rootCmd := root.NewRootCmd(cfg, logger, io, version, date)

	if err := rootCmd.Execute(); err != nil {
		logger.Error(fmt.Sprintf("execution error: %s\n", err))
	}
}
