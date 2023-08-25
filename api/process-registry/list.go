package registry

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *processRegistryClient) List(
	server *configuration.Server, productID, processType string,
) ([]ProcessRegistry, error) {
	query := `
		query ProcessRegistries($productID: ID!, $processType: String) {
			processRegistries(productID: $productID, processType: $processType) {
					id
					name
					version
					type
					image
					uploadDate
					owner
			}
		}
		`
	vars := map[string]interface{}{
		"productID":   productID,
		"processType": processType,
	}

	var respData struct {
		ProcessRegistries []ProcessRegistry `json:"processRegistries"`
	}

	err := c.client.MakeRequest(server, query, vars, &respData)

	return respData.ProcessRegistries, err
}
