package processregistry

import (
	"os"
	"time"

	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type RegisteredProcess struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	Type       string    `json:"type"`
	Image      string    `json:"image"`
	UploadDate time.Time `json:"uploadDate"`
	Owner      string    `json:"owner"`
	Status     string    `json:"status"`
}

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/registry_client.go -package=mocks

type APIClient interface {
	Register(server *configuration.Server, processFile *os.File, productID,
		processID, processType, version string) (*RegisteredProcess, error)
	List(server *configuration.Server, productID, processType string) ([]*RegisteredProcess, error)
}

type processRegistryClient struct {
	client *graphql.GqlManager
}

// NewClient creates a new struct to access Versions methods.
func NewClient(gql *graphql.GqlManager) APIClient {
	return &processRegistryClient{
		gql,
	}
}
