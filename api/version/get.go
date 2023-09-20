package version

import "github.com/konstellation-io/kli/internal/services/configuration"

func (c *Client) Get(server *configuration.Server, productID, versionTag string) (*Version, error) {
	query := `
		query Version($productID: ID!) {
			version(productID: $productID) {
					tag
					creationDate
					status
			}
		}
		`
	vars := map[string]interface{}{
		"productID":  productID,
		"versionTag": versionTag,
	}

	var respData struct {
		Version *Version `json:"version"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.Version, err
}
