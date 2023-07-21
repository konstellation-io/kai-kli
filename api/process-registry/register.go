package registry

import (
	"os"

	"github.com/konstellation-io/graphql"

	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *processRegistryClient) Register(server *configuration.Server, processFile *os.File,
	productID, processID, processType, version string) (string, error) {
	query := `
		mutation RegisterProcess($input: RegisterProcessInput!) {
			registerProcess(input: $input) {
				processedImageID
			}
		}
		`
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"file":        nil,
			"version":     version,
			"productID":   productID,
			"processID":   processID,
			"processType": processType,
		},
	}

	var respData struct {
		RegisteredProcess struct {
			ProcessedImageID string `json:"processedImageID"`
		} `json:"createdProcessImage"`
	}

	file := graphql.File{
		Field: "variables.input.file",
		Name:  processFile.Name(),
		R:     processFile,
	}

	err := c.client.UploadFile(server, file, query, vars, &respData)
	registeredProcessID := respData.RegisteredProcess.ProcessedImageID

	return registeredProcessID, err
}
