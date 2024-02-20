package processregistry

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *processRegistryClient) List(
	server *configuration.Server, productID, processName, version, processType string,
) ([]*entity.RegisteredProcess, error) {
	query := `
		query RegisteredProcesses($productID: ID!, $processName: String, $version: String, $processType: String) {
			registeredProcesses(productID: $productID, processName: $processName, version: $version, processType: $processType) {
					id
					name
					version
					type
					image
					uploadDate
					owner
					status
					isPublic
			}
		}
		`
	vars := map[string]interface{}{
		"productID":   productID,
		"processName": processName,
		"version":     version,
		"processType": processType,
	}

	var respData struct {
		RegisteredProcesses []*entity.RegisteredProcess `json:"registeredProcesses"`
	}

	err := c.client.MakeRequest(server, query, vars, &respData)

	return respData.RegisteredProcesses, err
}
