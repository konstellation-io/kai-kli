package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/api/kai/version"
	"github.com/konstellation-io/kli/cmd/args"
	config2 "github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/kai"
	"github.com/konstellation-io/kli/internal/logging"
)

// NewUpdateConfigCmd manage config command for version.
func NewUpdateConfigCmd(logger logging.Interface, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <version-name>",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Set config variables values",
		RunE: func(cmd *cobra.Command, args []string) error {
			versionName := args[0]

			product, _ := cmd.Flags().GetString("product")

			// read key=value pairs
			keyValuePairs, err := cmd.Flags().GetStringSlice("set")
			if err != nil {
				return err
			}

			// read values from env
			envVars, err := cmd.Flags().GetStringSlice("set-from-env")
			if err != nil {
				return err
			}
			if len(envVars) > 0 {
				keyValuePairs = addEnvVars(keyValuePairs, envVars)
			}

			// read values from file
			envFiles, err := cmd.Flags().GetStringSlice("set-from-file")
			if err != nil {
				return err
			}
			if len(envFiles) > 0 {
				keyValuePairs, err = addEnvFromFiles(keyValuePairs, envFiles)
				if err != nil {
					return err
				}
			}

			serverName, err := cmd.Flags().GetString("server")
			if err != nil {
				return err
			}

			server := cfg.GetByServerName(serverName)
			kaiClient := api.NewKaiClient(&graphql.ClientConfig{
				Debug:                 viper.GetBool(config2.DebugKey),
				DefaultRequestTimeout: viper.GetDuration("request_timeout"),
			}, server, viper.GetString(config2.BuildVersionKey))
			kaiInteractor := kai.NewInteractorWithDefaultRenderer(logger, kaiClient, cmd.OutOrStdout())

			if len(keyValuePairs) > 0 {
				logger.Info("No new configuration received")
			}

			configVars := getNewConfigVars(keyValuePairs)
			return kaiInteractor.UpdateVersionConfig(product, versionName, configVars)
		},
	}
	cmd.Flags().StringSlice("set", []string{}, "Set new key value pair key=value")
	cmd.Flags().StringSlice("set-from-env", []string{}, "Set new variable with value existing on current env")
	cmd.Flags().StringSlice("set-from-file", []string{}, "Set variables from a file with key/value pairs")

	return cmd
}

func addEnvFromFiles(pairs, envFiles []string) ([]string, error) {
	merged := pairs

	for _, file := range envFiles {
		fileVars, err := godotenv.Read(file)
		if err != nil {
			return nil, fmt.Errorf("error reading env file: %w", err)
		}

		for key, value := range fileVars {
			merged = append(merged, fmt.Sprintf("%s=%v", key, value))
		}
	}

	return merged, nil
}

func addEnvVars(pairs, envKeys []string) []string {
	merged := pairs

	for _, key := range envKeys {
		value := os.Getenv(key)
		merged = append(merged, fmt.Sprintf("%s=%v", key, value))
	}

	return merged
}

func getNewConfigVars(vars []string) []version.ConfigVariableInput {
	const substringNumber = 2

	newConfig := make([]version.ConfigVariableInput, 0, len(vars))

	for _, v := range vars {
		arr := strings.SplitN(v, "=", substringNumber)

		newConfig = append(newConfig, version.ConfigVariableInput{
			"key":   arr[0],
			"value": arr[1],
		})
	}

	return newConfig
}
