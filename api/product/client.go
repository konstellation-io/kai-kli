package product

import (
	"errors"

	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Client struct {
	client *graphql.GqlManager
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
				id
				name
			}
		}
	`

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"id":          "",
			"name":        name,
			"description": description,
		},
	}

	var respData struct {
		Products []kai.Product
	}

	err := c.client.MakeRequest(server, mutation, vars, &respData)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetProduct(server *configuration.Server, id string) (*kai.Product, error) {
	mutation := `
		query GetProduct($id: ID!) {
			product(id: $id) {
				id
				name
				description
				publishedVersion
			}
		}
	`

	vars := map[string]interface{}{
		"id": id,
	}

	var respData struct {
		Product *kai.Product
	}

	if err := c.client.MakeRequest(server, mutation, vars, &respData); err != nil {
		return nil, err
	}

	if respData.Product == nil {
		return nil, errors.New("no product returned")
	}

	return respData.Product, nil
}

func (c *Client) GetProducts(server *configuration.Server) ([]kai.Product, error) {
	mutation := `
		query GetProducts() {
			products() {
				id
				name
				description
			}
		}
	`

	vars := map[string]interface{}{}

	var respData struct {
		Products []kai.Product
	}

	if err := c.client.MakeRequest(server, mutation, vars, &respData); err != nil {
		return nil, err
	}

	products := make([]kai.Product, 0, len(respData.Products))

	for _, p := range respData.Products {
		products = append(products, kai.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	return products, nil
}
