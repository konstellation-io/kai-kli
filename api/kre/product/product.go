package product

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/${GOFILE} -package=mocks

import (
	"fmt"

	"github.com/konstellation-io/kli/api/graphql"
)

type Product struct {
	ID   string
	Name string
}

func GetProductsNames(products []Product) []string {
	names := make([]string, 0, len(products))
	for _, p := range products {
		names = append(names, p.Name)
	}

	return names
}

// ProductInterface method to interact with products.
type ProductInterface interface { //nolint:golint
	List() ([]Product, error)
	Create(name, description string) error
}

type productClient struct {
	client *graphql.GqlManager
}

// New creates a new struct to access Versions methods.
func New(gql *graphql.GqlManager) ProductInterface {
	return &productClient{
		gql,
	}
}

func (c *productClient) List() ([]Product, error) {
	query := `
		query GetProducts() {
			products() {
				id
				name
			}
		}
	`

	var respData struct {
		Products []Product
	}

	err := c.client.MakeRequest(query, nil, &respData)
	if err != nil {
		return nil, fmt.Errorf("product client error: %w", err)
	}

	return respData.Products, nil
}

func (c *productClient) Create(name, description string) error {
	mutation := `
		mutation CreateProduct($input: CreateProductInput!) {
			createProduct(input: $input) {
				name
			}
		}
	`

	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"id":          name,
			"name":        name,
			"description": description,
		},
	}

	var respData struct {
		Products []Product
	}

	err := c.client.MakeRequest(mutation, vars, &respData)
	if err != nil {
		return fmt.Errorf("product client error: %w", err)
	}

	return nil
}
