package processregistry

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *processRegistryClient) List(
	server *configuration.Server, productID, processType string,
) ([]*entity.RegisteredProcess, error) {
	query := `
		query RegisteredProcesses($productID: ID!, $processType: String) {
			registeredProcesses(productID: $productID, processType: $processType) {
					id
					name
					version
					type
					image
					uploadDate
					owner
					status
			}
		}
		`
	vars := map[string]interface{}{
		"productID":   productID,
		"processType": processType,
	}

	var respData struct {
		RegisteredProcesses []*entity.RegisteredProcess `json:"registeredProcesses"`
	}

	err := c.client.MakeRequest(server, query, vars, &respData)

	return respData.RegisteredProcesses, err
}
