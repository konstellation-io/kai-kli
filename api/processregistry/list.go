package processregistry

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *processRegistryClient) List(
	server *configuration.Server, productID, processType string,
) ([]RegisteredProcess, error) {
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
			}
		}
		`
	vars := map[string]interface{}{
		"productID":   productID,
		"processType": processType,
	}

	var respData struct {
		RegisteredProcesses []RegisteredProcess `json:"registeredProcesses"`
	}

	err := c.client.MakeRequest(server, query, vars, &respData)

	return respData.RegisteredProcesses, err
}
