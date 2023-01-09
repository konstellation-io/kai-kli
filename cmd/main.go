package main

import (
	"fmt"

	"github.com/konstellation-io/kli/cmd/factory"
	"github.com/konstellation-io/kli/pkg/cmd/root"
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

	cmdFactory := factory.NewCmdFactory(buildVersion)

	rootCmd := root.NewRootCmd(cmdFactory, buildVersion, buildDate)

	if err := rootCmd.Execute(); err != nil {
		cmdFactory.Logger().Error(fmt.Sprintf("execution error: %s\n", err))
	}
}
