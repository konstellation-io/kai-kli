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
	RenderLogs(productID string, logs []entity.Log, outFormat entity.LogOutFormat, showAllLabels bool)
	RenderCallout(v *entity.Version)
	RenderKliVersion(version, buildDate string)
	RenderProcessRegistered(process *entity.RegisteredProcess)
	RenderProcessDeleted(process string)
	RenderProductCreated(product string, server *configuration.Server, initLocal bool)
	RenderProductBinded(productID string)
	RenderLogin(serverName string)
	RenderLogout(serverName string)
	RenderPushVersion(product, versionTag string)
	RenderStartVersion(product, versionTag string)
	RenderStopVersion(product, versionTag string)
	RenderPublishVersion(product, versionTag string, triggers []entity.TriggerEndpoint)
	RenderUnpublishVersion(product, versionTag string)
	RenderAddUserToProduct(product, userEmail string)
	RenderRemoveUserFromProduct(product, userEmail string)
	RenderAddMaintainerToProduct(product, userEmail string)
	RenderRemoveMaintainerFromProduct(product, userEmail string)
}
