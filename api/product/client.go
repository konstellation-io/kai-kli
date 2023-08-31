package product

import (
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Client struct {
	client *graphql.GqlManager
}

type ProductDTO struct {
	Name string
}

func NewClient(gqlClient *graphql.GqlManager) *Client {
	return &Client{
		client: gqlClient,
	}
}

func (c *Client) CreateProduct(server *configuration.Server, name, description string) error {
	mutation := `
		mutation CreateProduct($input: CreateProductInput!) {
			createProduct(input: $input) {
				name
			}
		}
	`

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"name":        name,
			"description": description,
		},
	}

	var respData struct {
		Products []ProductDTO
	}

	return c.client.MakeRequest(server, mutation, vars, &respData)
}
