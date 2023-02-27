package main

import (
	"fmt"

	"github.com/konstellation-io/kli/api/kre/config"
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
	}

	rootCmd := root.NewRootCmd(cfg, logger, io, buildVersion, buildDate)

	if err := rootCmd.Execute(); err != nil {
		logger.Error(fmt.Sprintf("execution error: %s\n", err))
	}
}
