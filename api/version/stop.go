package version

import "github.com/konstellation-io/kli/internal/services/configuration"

func (c *Client) Stop(server *configuration.Server, productID, versionTag, comment string) (string, error) {
	query := `
		mutation StopVersion($input: StopVersionInput!) {
			stopVersion(input: $input) {
				tag
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
		StopVersion struct {
			Tag string `json:"tag"`
		} `json:"stopVersion"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.StopVersion.Tag, err
}
