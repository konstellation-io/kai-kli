package version

import (
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

type Handler struct {
	logger         logging.Interface
	renderer       render.Renderer
	versionClient  kai.VersionClient
	productClient  kai.ProductClient
	configService  *configuration.KaiConfigService
	productService *productconfiguration.ProductConfigService
}

func NewHandler(
	logger logging.Interface,
	renderer render.Renderer,
	versionClient kai.VersionClient,
	productClient kai.ProductClient,
) *Handler {
	return &Handler{
		logger:         logger,
		renderer:       renderer,
		versionClient:  versionClient,
		productClient:  productClient,
		configService:  configuration.NewKaiConfigService(logger),
		productService: productconfiguration.NewProductConfigService(logger),
	}
}
