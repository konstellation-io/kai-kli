package render

import (
	"github.com/konstellation-io/kli/internal/services/configuration"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/renderer.go -package=mocks

// Renderer interface that can be used to render in different formats.
type Renderer interface {
	RenderServers(servers []configuration.Server)
}
