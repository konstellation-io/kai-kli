package render

import (
	"github.com/konstellation-io/kli/api/kai/config"
	"github.com/konstellation-io/kli/api/kai/product"
	"github.com/konstellation-io/kli/api/kai/version"
	"github.com/konstellation-io/kli/internal/krt/errors"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/renderer.go -package=mocks

// Renderer interface that can be used to render in different formats.
type Renderer interface {
	RenderServerList(servers []config.ServerConfig, defaultServer string)
	RenderVars(cfg *version.Config, showValues bool)
	RenderVersions(versions version.List)
	RenderValidationErrors(validationErrors []*errors.ValidationError)
	RenderProducts(products []product.Product)
}
