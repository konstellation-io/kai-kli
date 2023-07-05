package version

func (c *versionClient) Publish(product, versionName, comment string) error {
	query := `
		mutation PublishVersion($input: PublishVersionInput!) {
			publishVersion(input: $input) {
				status
			}
		}
	`
	vars := map[string]interface{}{
		"input": map[string]string{
			"productID":   product,
			"versionName": versionName,
			"comment":     comment,
		},
	}

	var respData struct {
		Status string
	}

	return c.client.MakeRequest(query, vars, &respData)
}
