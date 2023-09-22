package processregistry

import (
	"fmt"
	"os"

	"github.com/konstellation-io/graphql"
	"github.com/konstellation-io/kli/internal/entity"

	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *processRegistryClient) Register(server *configuration.Server, processFile *os.File,
	productID, processID, processType, version string) (*entity.RegisteredProcess, error) {
	query := `
		mutation RegisterProcess($input: RegisterProcessInput!) {
			registerProcess(input: $input) {
				id
				image
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
		RegisteredProcess entity.RegisteredProcess `json:"registerProcess"`
	}

	processFileOpen, err := os.Open(processFile.Name())
	if err != nil {
		return nil, fmt.Errorf("reading file content: %w", err)
	}

	defer processFileOpen.Close()

	file := graphql.File{
		Field: "variables.input.file",
		Name:  processFileOpen.Name(),
		R:     processFileOpen,
	}

	err = c.client.UploadFile(server, file, query, vars, &respData)

	return &respData.RegisteredProcess, err
}
