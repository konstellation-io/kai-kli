package version

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Client) Publish(server *configuration.Server, productID, versionTag, comment string) ([]entity.Trigger, error) {
	query := `
		mutation PublishVersion($input: PublishVersionInput!) {
			publishVersion(input: $input) {
				trigger
				url
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
		PublishedTriggers []entity.Trigger `json:"publishVersion"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.PublishedTriggers, err
}
