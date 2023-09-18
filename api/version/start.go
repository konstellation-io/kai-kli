package version

import "github.com/konstellation-io/kli/internal/services/configuration"

func (c *Client) Start(server *configuration.Server, productID, versionTag, comment string) (string, string, error) {
	query := `
		mutation StartVersion($input: StartVersionInput!) {
			startVersion(input: $input) {
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
		StartVersion struct {
			Tag    string `json:"tag"`
			Status string `json:"status"`
		} `json:"startVersion"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.StartVersion.Tag, respData.StartVersion.Status, err
}
