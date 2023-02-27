package kre_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/kre/runtime"
)

func TestRuntimeList(t *testing.T) {
	srv, client := gqlMockServer(t, "", `
		{
				"data": {
						"runtimes": [
								{
										"name": "test-runtime"
								}
						]
				}
		}
	`)
	defer srv.Close()

	c := runtime.New(client)

	list, err := c.List()
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, list[0], runtime.Runtime{
		Name: "test-runtime",
	})
}

func TestRuntimeCreate(t *testing.T) {
	srv, client := gqlMockServer(t, "", `
		{
				"data": {
						"createRuntime": {
								"name": "test-runtime"
						}
				}
		}
	`)
	defer srv.Close()

	c := runtime.New(client)

	err := c.Create("test-runtime", "test description")
	require.NoError(t, err)
}
