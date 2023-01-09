package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/api/kre/version"
	"github.com/konstellation-io/kli/cmd/factory"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/pkg/cmd/args"
)

// NewConfigCmd manage config command for version.
func NewConfigCmd(f factory.CmdFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: "Get or set config values",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, _ := cmd.Flags().GetString("server")

			versionName := args[0]

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

			// Update config
			if len(keyValuePairs) > 0 {
				return updateConfig(f, serverName, versionName, keyValuePairs)
			}

			return getConfig(f, cmd, serverName, versionName)
		},
	}
	cmd.Flags().StringSlice("set", []string{}, "Set new key value pair key=value")
	cmd.Flags().StringSlice("set-from-env", []string{}, "Set new variable with value existing on current env")
	cmd.Flags().StringSlice("set-from-file", []string{}, "Set variables from a file with key/value pairs")
	cmd.Flags().Bool("show-values", false, "Show configuration variables")

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

func getConfig(f factory.CmdFactory, cmd *cobra.Command, serverName, versionName string) error {
	log := f.Logger()

	c, err := f.KreClient(serverName)
	if err != nil {
		return err
	}

	config, err := c.Version().GetConfig(versionName)
	if err != nil {
		return err
	}

	if len(config.Vars) == 0 {
		log.Info("No config found.")
		return nil
	}

	show, err := cmd.Flags().GetBool("show-values")
	if err != nil {
		return err
	}

	r := render.DefaultRenderer(cmd.OutOrStdout())
	renderVariables(r, config, show)

	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	if config.Completed {
		log.Success("Version config complete")
	} else {
		log.Warn("Version config incomplete")
	}

	return nil
}

func updateConfig(f factory.CmdFactory, serverName, versionName string, vars []string) error {
	log := f.Logger()

	const substringNumber = 2

	c, err := f.KreClient(serverName)
	if err != nil {
		return err
	}

	config := []version.ConfigVariableInput{}

	for _, v := range vars {
		arr := strings.SplitN(v, "=", substringNumber)

		config = append(config, version.ConfigVariableInput{
			"key":   arr[0],
			"value": arr[1],
		})
	}

	completed, err := c.Version().UpdateConfig(versionName, config)
	if err != nil {
		return err
	}

	status := "updated"
	if completed {
		status = "completed"
	}

	log.Success(fmt.Sprintf("Config %s for version '%s'.", status, versionName))

	return nil
}

func renderVariables(r render.Renderer, config *version.Config, showValues bool) {
	h := []string{
		"",
		"TYPE",
		"KEY",
	}

	if showValues {
		h = append(h, "VALUE")
	}

	r.SetHeader(h)

	for i, v := range config.Vars {
		s := []string{
			fmt.Sprint(i + 1),
			string(v.Type),
			v.Key,
		}

		if showValues {
			s = append(s, v.Value)
		}

		r.Append(s)
	}

	r.Render()
}
