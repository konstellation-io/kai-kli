package version

import "github.com/konstellation-io/kli/internal/services/configuration"

func (c *Client) Stop(server *configuration.Server, productID, versionTag, comment string) (string, error) {
	query := `
		mutation StopVersion($input: StopVersionInput!) {
			stopVersion(input: $input) {
				id
				tag
				status
			}
		}
		`
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"productID":  productID,
			"versionTag": versionTag,
			"comment":    comment,
		},
	}

	var respData struct {
		Tag string `json:"tag"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.Tag, err
}
