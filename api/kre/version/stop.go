package version

func (c *versionClient) Stop(runtime, versionName, comment string) error {
	query := `
		mutation StopVersion($input: StopVersionInput!) {
			stopVersion(input: $input) {
				status
			}
		}
	`
	vars := map[string]interface{}{
		"input": map[string]string{
			"runtimeId":   runtime,
			"versionName": versionName,
			"comment":     comment,
		},
	}

	var respData struct {
		Status string
	}

	return c.client.MakeRequest(query, vars, &respData)
}
