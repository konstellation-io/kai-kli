package api

import (
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/api/kre/runtime"
	"github.com/konstellation-io/kli/api/kre/version"
)

// KRE object to implement access to KRE API.
type KRE struct {
	gqlManager *graphql.GqlManager
	version    version.VersionInterface
	runtime    runtime.RuntimeInterface
}

// Version access to methods to interact with Versions.
func (a *KRE) Version() version.VersionInterface {
	return a.version
}

func (a *KRE) Runtime() runtime.RuntimeInterface {
	return a.runtime
}

// NewKreClient creates an API client instance.
func NewKreClient(cfg *graphql.ClientConfig, server *config.ServerConfig, appVersion string) *KRE {
	g := graphql.NewGqlManager(cfg, server, appVersion)

	return &KRE{
		g,
		version.New(g),
		runtime.New(g),
	}
}
