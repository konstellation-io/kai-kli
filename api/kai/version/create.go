package version

import (
	"os"

	"github.com/konstellation-io/graphql"
)

func (c *versionClient) Create(product, krtFile string) (string, error) {
	query := `
		mutation CreateVersion($input: CreateVersionInput!) {
			createVersion(input: $input) {
				name
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
			Name string `json:"name"`
		} `json:"createVersion"`
	}

	r, err := os.Open(krtFile)
	if err != nil {
		return "", err
	}

	file := graphql.File{
		Field: "variables.input.file",
		Name:  r.Name(),
		R:     r,
	}

	err = c.client.UploadFile(file, query, vars, &respData)
	versionName := respData.CreateVersion.Name

	return versionName, err
}
