package product

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Client) AddMaintainerToProduct(server *configuration.Server, product, userEmail string) error {
	mutation := `
		mutation AddMaintainerToProduct($input: AddUserToProductInput!) {
			addMaintainerToProduct(input: $input) {
				id
			}
		}
	`

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"product": product,
			"email":   userEmail,
		},
	}

	var respData struct {
		User *struct{ id string }
	}

	err := c.client.MakeRequest(server, mutation, vars, &respData)

	if err != nil {
		return err
	}

	return nil
}
