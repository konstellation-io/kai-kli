package version

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Client) List(server *configuration.Server, productID string, status *string) ([]*entity.Version, error) {
	query := `
		query Versions($productID: ID!, $status: String) {
			versions(productID: $productID, status: $status) {
					tag
					creationDate
					status
			}
		}
		`
	vars := map[string]interface{}{
		"productID": productID,
	}

	if status != nil {
		vars["status"] = *status
	}

	var respData struct {
		Versions []*entity.Version `json:"versions"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.Versions, err
}
