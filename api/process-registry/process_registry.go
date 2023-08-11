package registry

import (
	"os"

	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/registry_client.go -package=mocks

type ProcessRegistryInterface interface {
	Register(server *configuration.Server, processFile *os.File, productID,
		processID, processType, version string) (string, error)
}

type processRegistryClient struct {
	client *graphql.GqlManager
}

// New creates a new struct to access Versions methods.
func New(gql *graphql.GqlManager) ProcessRegistryInterface {
	return &processRegistryClient{
		gql,
	}
}
