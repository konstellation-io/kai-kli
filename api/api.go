package api

import (
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/processregistry"
)

// KAI object to implement access to KAI API.
type KAI struct {
	gqlManager        *graphql.GqlManager
	registeredProcess processregistry.ProcessRegistryInterface
}

// RegisteredProcess access to methods to interact with the Process Registry.
func (a *KAI) RegisteredProcess() processregistry.ProcessRegistryInterface {
	return a.registeredProcess
}

// NewKaiClient creates an API client instance.
func NewKaiClient() *KAI {
	g := graphql.NewGqlManager()

	return &KAI{
		g,
		processregistry.New(g),
	}
}
