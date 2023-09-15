package version

import (
	"fmt"
	"os"

	"github.com/konstellation-io/graphql"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Client) Push(server *configuration.Server, product, krtFilePath string) (string, error) {
	query := `
		mutation CreateVersion($input: CreateVersionInput!) {
			createVersion(input: $input) {
				tag
			}
		}
		`

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"file":      nil,
			"productID": product,
		},
	}

	var respData struct {
		CreateVersion struct {
			Tag string `json:"tag"`
		} `json:"createVersion"`
	}

	krtFile, err := os.Open(krtFilePath)
	if err != nil {
		return "", fmt.Errorf("opening krt file: %w", err)
	}

	file := graphql.File{
		Field: "variables.input.file",
		Name:  krtFile.Name(),
		R:     krtFile,
	}

	err = c.gqlClient.UploadFile(server, file, query, vars, &respData)
	if err != nil {
		return "", err
	}

	return respData.CreateVersion.Tag, nil
}
