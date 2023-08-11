package kai

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/kaiclient.go -package=mocks

import (
	processregistry "github.com/konstellation-io/kli/api/process-registry"
)

// Interface first level methods.
type Interface interface { //nolint: golint
	ProcessRegistry() processregistry.ProcessRegistryInterface
}
