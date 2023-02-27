package runtime

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/${GOFILE} -package=mocks

import (
	"fmt"

	"github.com/konstellation-io/kli/api/graphql"
)

type Runtime struct {
	Name string
}

func GetRuntimesNames(runtimes []Runtime) []string {
	names := make([]string, 0, len(runtimes))
	for _, r := range runtimes {
		names = append(names, r.Name)
	}

	return names
}

// RuntimeInterface method to interact with Runtimes.
type RuntimeInterface interface { //nolint:golint
	List() ([]Runtime, error)
	Create(name, description string) error
}

type runtimeClient struct {
	client *graphql.GqlManager
}

// New creates a new struct to access Versions methods.
func New(gql *graphql.GqlManager) RuntimeInterface {
	return &runtimeClient{
		gql,
	}
}

func (c *runtimeClient) List() ([]Runtime, error) {
	query := `
		query GetRuntimes() {
			runtimes() {
				name
			}
		}
	`

	var respData struct {
		Runtimes []Runtime
	}

	err := c.client.MakeRequest(query, nil, &respData)
	if err != nil {
		return nil, fmt.Errorf("runtime client error: %w", err)
	}

	return respData.Runtimes, nil
}

func (c *runtimeClient) Create(name, description string) error {
	mutation := `
		mutation CreateRuntime($input: CreateRuntimeInput!) {
			createRuntime(input: $input) {
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
		Runtimes []Runtime
	}

	err := c.client.MakeRequest(mutation, vars, &respData)
	if err != nil {
		return fmt.Errorf("runtime client error: %w", err)
	}

	return nil
}
