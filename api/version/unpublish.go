package version

import "github.com/konstellation-io/kli/internal/services/configuration"

func (c *Client) Unpublish(server *configuration.Server, productID, versionTag, comment string) (string, error) {
	query := `
		mutation UnpublishVersion($input: UnpublishVersionInput!) {
			unpublishVersion(input: $input) {
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
		UnpublishVersion struct {
			Tag string `json:"tag"`
		} `json:"unpublishVersion"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.UnpublishVersion.Tag, err
}
