package kai

import (
	"github.com/konstellation-io/kli/api/processregistry"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

//go:generate mockery --name Client --output ../../mocks --filename kaiclient.go --structname MockKaiClient
//go:generate mockery --name ProductClient --output ../../mocks --filename product_client.go --structname MockProductClient

// Client first level methods.
type Client interface { //nolint: golint
	ProcessRegistry() processregistry.APIClient
	ProductClient() ProductClient
}

type ProductClient interface {
	CreateProduct(server *configuration.Server, name, description string) error
}
