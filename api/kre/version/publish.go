package version

func (c *versionClient) Publish(runtime, versionName, comment string) error {
	query := `
		mutation PublishVersion($input: PublishVersionInput!) {
			publishVersion(input: $input) {
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
