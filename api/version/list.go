package version

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Client) List(server *configuration.Server, productID string) ([]*entity.Version, error) {
	query := `
		query Versions($productID: ID!) {
			versions(productID: $productID) {
					tag
					creationDate
					status
			}
		}
		`
	vars := map[string]interface{}{
		"productID": productID,
	}

	var respData struct {
		Versions []*entity.Version `json:"versions"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.Versions, err
}
