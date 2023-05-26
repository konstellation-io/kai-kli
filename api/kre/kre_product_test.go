package kre_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/kre/product"
)

func TestProductList(t *testing.T) {
	srv, client := gqlMockServer(t, "", `
		{
				"data": {
						"products": [
								{
										"id": "test-product",
										"name": "Test product"
								}
						]
				}
		}
	`)
	defer srv.Close()

	c := product.New(client)

	list, err := c.List()
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, list[0], product.Product{
		ID:   "test-product",
		Name: "Test product",
	})
}

func TestProductCreate(t *testing.T) {
	srv, client := gqlMockServer(t, "", `
		{
				"data": {
						"createProduct": {
								"name": "test-product"
						}
				}
		}
	`)
	defer srv.Close()

	c := product.New(client)

	err := c.Create("test-product", "test description")
	require.NoError(t, err)
}
