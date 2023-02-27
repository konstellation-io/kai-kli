package version

func (c *versionClient) Start(runtime, versionName, comment string) error {
	query := `
		mutation StartVersion($input: StartVersionInput!) {
			startVersion(input: $input) {
				status
			}
		}
	`
	vars := map[string]interface{}{
		"input": map[string]string{
			"versionName": versionName,
			"comment":     comment,
			"runtimeId":   runtime,
		},
	}

	var respData struct {
		Status string
	}

	return c.client.MakeRequest(query, vars, &respData)
}
