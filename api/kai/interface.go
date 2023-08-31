package kai

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/kaiclient.go -package=mocks

import (
	processregistry "github.com/konstellation-io/kli/api/process-registry"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

// Interface first level methods.
type Interface interface { //nolint: golint
	ProcessRegistry() processregistry.ProcessRegistryInterface
	ProductClient() ProductClient
}

type ProductClient interface {
	CreateProduct(server *configuration.Server, name, description string) error
}
