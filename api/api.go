package api

import (
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/api/kai/product"
	"github.com/konstellation-io/kli/api/kai/version"
)

// KAI object to implement access to KAI API.
type KAI struct {
	gqlManager *graphql.GqlManager
	version    version.VersionInterface
	product    product.ProductInterface
}

// Version access to methods to interact with Versions.
func (a *KAI) Version() version.VersionInterface {
	return a.version
}

func (a *KAI) Product() product.ProductInterface {
	return a.product
}

// NewKAIClient creates an API client instance.
func NewKAIClient(cfg *graphql.ClientConfig, server *config.ServerConfig, appVersion string) *KAI {
	g := graphql.NewGqlManager(cfg, server, appVersion)

	return &KAI{
		g,
		version.New(g),
		product.New(g),
	}
}
