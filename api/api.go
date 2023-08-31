package api

import (
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai"
	registry "github.com/konstellation-io/kli/api/process-registry"
	"github.com/konstellation-io/kli/api/product"
)

// KAI object to implement access to KAI API.
type KAI struct {
	gqlManager      *graphql.GqlManager
	processRegistry registry.ProcessRegistryInterface
	productClient   kai.ProductClient
}

// ProcessRegistry access to methods to interact with the Process Registry.
func (a *KAI) ProcessRegistry() registry.ProcessRegistryInterface {
	return a.processRegistry
}

// ProductClient access to methods to interact with the Process Registry.
func (a *KAI) ProductClient() kai.ProductClient {
	return a.productClient
}

// NewKaiClient creates an API client instance.
func NewKaiClient() *KAI {
	g := graphql.NewGqlManager()

	return &KAI{
		g,
		registry.New(g),
		product.NewClient(g),
	}
}
