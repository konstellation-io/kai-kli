package product

import (
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

type Handler struct {
	logger        logging.Interface
	renderer      render.Renderer
	configService *productconfiguration.ProductConfigService
	productClient kai.ProductClient
	versionClient kai.VersionClient
}

func NewHandler(
	logger logging.Interface,
	renderer render.Renderer,
	productClient kai.ProductClient,
	versionClient kai.VersionClient,
	configService *productconfiguration.ProductConfigService,
) *Handler {
	return &Handler{
		logger:        logger,
		renderer:      renderer,
		configService: configService,
		productClient: productClient,
		versionClient: versionClient,
	}
}
