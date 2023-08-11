package api

import (
	"github.com/konstellation-io/kli/api/graphql"
	registry "github.com/konstellation-io/kli/api/process-registry"
)

// KAI object to implement access to KAI API.
type KAI struct {
	gqlManager      *graphql.GqlManager
	processRegistry registry.ProcessRegistryInterface
}

// ProcessRegistry access to methods to interact with the Process Registry.
func (a *KAI) ProcessRegistry() registry.ProcessRegistryInterface {
	return a.processRegistry
}

// NewKaiClient creates an API client instance.
func NewKaiClient() *KAI {
	g := graphql.NewGqlManager()

	return &KAI{
		g,
		registry.New(g),
	}
}
