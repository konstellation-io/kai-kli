package version

import "github.com/konstellation-io/kli/internal/services/configuration"

func (c *Client) Get(server *configuration.Server, productID, versionTag string) (*Version, error) {
	query := `
		query Version($productID: ID!, $tag: String!) {
			version(productID: $productID, tag: $tag) {
					tag
					creationDate
					status
					error
			}
		}
		`
	vars := map[string]interface{}{
		"productID": productID,
		"tag":       versionTag,
	}

	var respData struct {
		Version *Version `json:"version"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.Version, err
}
