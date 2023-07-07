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

const (
	// Version is dynamically set by the toolchain or overridden at build time.
	Version = "DEV"

	// Date is dynamically set at build time.
	Date = "" // YYYY-MM-DD
)

func main() {
	buildDate := Date
	buildVersion := Version

	logger := logging.NewDefaultLogger()
	io := iostreams.System()

	cfg, err := config.NewConfig(buildVersion)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = config2.InitConfigWithBuildVersion(buildVersion)
	if err != nil {
		log.Fatalf("initialize configuration: %s", err)
	}

	rootCmd := root.NewRootCmd(cfg, logger, io, buildVersion, buildDate)

	if err := rootCmd.Execute(); err != nil {
		logger.Error(fmt.Sprintf("execution error: %s\n", err))
	}
}
