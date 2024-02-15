package processregistry

import (
	"os"

	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/registry_client.go -package=mocks

type APIClient interface {
	Register(server *configuration.Server, processFile *os.File, productID,
		processID, processType, version string) (*entity.RegisteredProcess, error)
	RegisterPublic(server *configuration.Server, processFile *os.File, processID, processType,
		version string) (*entity.RegisteredProcess, error)
	List(server *configuration.Server, productID, processType string) ([]*entity.RegisteredProcess, error)
	Delete(server *configuration.Server, productID, processID, version string) (string, error)
	DeletePublic(server *configuration.Server, processID, version string) (string, error)
}

type processRegistryClient struct {
	client *graphql.GqlManager
}

// NewClient creates a new struct to access registry methods.
func NewClient(gql *graphql.GqlManager) APIClient {
	return &processRegistryClient{
		gql,
	}
}
