package kai

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/kaiclient.go -package=mocks

import (
	"github.com/konstellation-io/kli/api/processregistry"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

// Interface first level methods.
type Interface interface { //nolint: golint
	ProcessRegistry() processregistry.APIClient
	ProductClient() ProductClient
}

type ProductClient interface {
	CreateProduct(server *configuration.Server, name, description string) error
}
