package kai

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/kaiclient.go -package=mocks

import (
	"github.com/konstellation-io/kli/api/processregistry"
)

// Interface first level methods.
type Interface interface { //nolint: golint
	ProcessRegistry() processregistry.APIClient
}
