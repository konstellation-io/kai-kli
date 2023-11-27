package render

import (
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/krt/pkg/krt"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/renderer.go -package=mocks

// Renderer interface that can be used to render in different formats.
type Renderer interface {
	RenderServers(servers []*configuration.Server)
	RenderWorkflows(workflows []krt.Workflow)
	RenderProcesses(processes []krt.Process)
	RenderConfiguration(scope string, config map[string]string)
	RenderRegisteredProcesses(registries []*entity.RegisteredProcess)
	RenderProducts(products []kai.Product)
	RenderVersion(productID string, v *entity.Version)
	RenderVersions(productID string, versions []*entity.Version)
	RenderTriggers(triggers []entity.TriggerEndpoint)
	RenderLogs(logs []entity.Log, outFormat entity.LogOutFormat, showAllLabels bool)
}
