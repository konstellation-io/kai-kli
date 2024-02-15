package processregistry

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *processRegistryClient) Delete(
	server *configuration.Server,
	productID, processID, version string,
) (string, error) {
	query := `
	mutation DeleteProcess($input: DeleteProcessInput!) {
		deleteProcess(input: $input)
	}
`

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"productID": productID,
			"processID": processID,
			"version":   version,
		},
	}

	var respData struct {
		DeletedProcessID string `json:"deleteProcess"`
	}

	err := c.client.MakeRequest(server, query, vars, &respData)

	return respData.DeletedProcessID, err
}

func (c *processRegistryClient) DeletePublic(
	server *configuration.Server,
	processID, version string,
) (string, error) {
	query := `
	mutation DeletePublicProcess($input: DeletePublicProcessInput!) {
		deletePublicProcess(input: $input)
	}
`

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"processID": processID,
			"version":   version,
		},
	}

	var respData struct {
		DeletedProcessID string `json:"deletePublicProcess"`
	}

	err := c.client.MakeRequest(server, query, vars, &respData)

	return respData.DeletedProcessID, err
}
