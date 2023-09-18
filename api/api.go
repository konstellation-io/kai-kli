package api

import (
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/api/processregistry"
	"github.com/konstellation-io/kli/api/product"
	"github.com/konstellation-io/kli/api/version"
)

// KAI object to implement access to KAI API.
type KAI struct {
	gqlManager      *graphql.GqlManager
	processRegistry processregistry.APIClient
	productClient   kai.ProductClient
	versionClient   kai.VersionClient
}

// ProcessRegistry access to methods to interact with the Process Registry.
func (a *KAI) ProcessRegistry() processregistry.APIClient {
	return a.processRegistry
}

// ProductClient access to methods to interact with the product's API.
func (a *KAI) ProductClient() kai.ProductClient {
	return a.productClient
}

// VersionClient access to methods to interact with the version's API.
func (a *KAI) VersionClient() kai.VersionClient {
	return a.versionClient
}

// NewKaiClient creates an API client instance.
func NewKaiClient() *KAI {
	g := graphql.NewGqlManager()

	return &KAI{
		g,
		processregistry.NewClient(g),
		product.NewClient(g),
		version.NewClient(g),
	}
}
