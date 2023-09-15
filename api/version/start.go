package version

import "github.com/konstellation-io/kli/internal/services/configuration"

func (c *Client) Start(server *configuration.Server, productID, versionTag, comment string) (*Version, error) {
	query := `
		mutation StartVersion($input: StartVersionInput!) {
			startVersion(input: $input) {
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
		StartVersion Version `json:"startVersion"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return &respData.StartVersion, err
}
