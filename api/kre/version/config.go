package version

// Config struct to represent a Version config.
type Config struct {
	Completed bool
	Vars      []*ConfigVariable
}

// ConfigVariable struct to represent a Version config variables.
type ConfigVariable struct {
	Key   string
	Value string
	Type  ConfigVariableType
}

// ConfigVariableInput struct to collect key/values to update Version config.
type ConfigVariableInput map[string]string

// ConfigVariableType enum to represent variable types.
type ConfigVariableType string

const (
	// ConfigVariableTypeVariable type VARIABLE.
	ConfigVariableTypeVariable ConfigVariableType = "VARIABLE"
	// ConfigVariableTypeFile type FILE.
	ConfigVariableTypeFile ConfigVariableType = "FILE"
)

// IsValid method to implement enums.
func (e ConfigVariableType) IsValid() bool {
	switch e {
	case ConfigVariableTypeVariable, ConfigVariableTypeFile:
		return true
	}

	return false
}

func (e ConfigVariableType) String() string {
	return string(e)
}

// UpdateConfig update a Version config values.
func (c *versionClient) UpdateConfig(runtime, versionName string, configVars []ConfigVariableInput) (bool, error) {
	query := `
		mutation UpdateConfig($input: UpdateConfigurationInput!) {
			updateVersionUserConfiguration(input: $input) {
				config {
					completed
				}
			}
		}
	`
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"versionName":            versionName,
			"configurationVariables": configVars,
			"runtimeId":              runtime,
		},
	}

	var respData struct {
		UpdateVersionUserConfiguration struct {
			Config struct {
				Completed bool
			}
		}
	}

	err := c.client.MakeRequest(query, vars, &respData)

	return respData.UpdateVersionUserConfiguration.Config.Completed, err
}

func (c *versionClient) GetConfig(runtime, versionName string) (*Config, error) {
	query := `
    query GetVersionConf($versionName: String!, $runtimeId: ID!) {
      version(name: $versionName, runtimeId: $runtimeId) {
        config {
          completed
          vars{
            key
            value
            type
          }
        }
      }
    }
  `
	vars := map[string]interface{}{
		"versionName": versionName,
		"runtimeId":   runtime,
	}

	var respData struct {
		Version struct {
			Config Config
		}
	}

	err := c.client.MakeRequest(query, vars, &respData)

	return &respData.Version.Config, err
}
