package render

import (
	"github.com/konstellation-io/krt/pkg/krt"

	registry "github.com/konstellation-io/kli/api/process-registry"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/renderer.go -package=mocks

// Renderer interface that can be used to render in different formats.
type Renderer interface {
	RenderServers(servers []configuration.Server)
	RenderWorkflows(workflows []krt.Workflow)
	RenderProcesses(processes []krt.Process)
	RenderConfiguration(scope string, config map[string]string)
	RenderProcessRegistries(registries []registry.ProcessRegistry)
}
