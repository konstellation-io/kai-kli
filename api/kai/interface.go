package kai

import (
	"github.com/konstellation-io/kli/api/processregistry"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

//go:generate mockery --name Client --output ../../mocks --filename kaiclient.go --structname MockKaiClient
//go:generate mockery --name ProductClient --output ../../mocks --filename product_client.go --structname MockProductClient
//go:generate mockery --name VersionClient --output ../../mocks --filename version_client.go --structname MockVersionClient

type Product struct {
	ID          string
	Name        string
	Description string
}

// Client first level methods.
type Client interface { //nolint: golint
	ProcessRegistry() processregistry.APIClient
	ProductClient() ProductClient
}

type ProductClient interface {
	CreateProduct(server *configuration.Server, name, description string) error
	GetProduct(server *configuration.Server, id string) (*Product, error)
	GetProducts(server *configuration.Server) ([]Product, error)
}

type VersionClient interface {
	Push(server *configuration.Server, product, krtFilePath string) (startedVersionTag string, err error)
	Start(server *configuration.Server, productID, versionTag, comment string) (tag string, status string, err error)
}
