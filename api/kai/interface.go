package kai

import (
	"github.com/konstellation-io/kli/api/processregistry"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

//go:generate mockery --name Client --output ../../mocks --filename kaiclient.go --structname MockKaiClient
//go:generate mockery --name ProductClient --output ../../mocks --filename product_client.go --structname MockProductClient
//go:generate mockery --name VersionClient --output ../../mocks --filename version_client.go --structname MockVersionClient

type Product struct {
	ID               string
	Name             string
	Description      string
	PublishedVersion *string
}

// Client first level methods.
type Client interface { //nolint: golint
	ProcessRegistry() processregistry.APIClient
	ProductClient() ProductClient
}

type ProductClient interface {
	CreateProduct(server *configuration.Server, name, description string) error
	GetProduct(server *configuration.Server, id string) (*Product, error)
	GetProducts(server *configuration.Server, productName string) ([]Product, error)
	AddUserToProduct(server *configuration.Server, product, userEmail string) error
	RemoveUserFromProduct(server *configuration.Server, product, userEmail string) error
	AddMaintainerToProduct(server *configuration.Server, product, userEmail string) error
	RemoveMaintainerFromProduct(server *configuration.Server, product, userEmail string) error
}

type VersionClient interface {
	Push(server *configuration.Server, product, krtFilePath string) (string, error)
	Start(server *configuration.Server, productID, versionTag, comment string) (string, error)
	Stop(server *configuration.Server, productID, versionTag, comment string) (string, error)
	Get(server *configuration.Server, productID string, versionTag *string) (*entity.Version, error)
	List(server *configuration.Server, productID string, status *string) ([]*entity.Version, error)
	Publish(server *configuration.Server, productID, versionTag, comment string, force bool) ([]entity.TriggerEndpoint, error)
	Unpublish(server *configuration.Server, productID, versionTag, comment string) (string, error)
	Logs(server *configuration.Server, filters *entity.LogFilters) ([]entity.Log, error)
}
