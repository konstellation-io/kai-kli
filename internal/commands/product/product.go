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
}

func NewHandler(logger logging.Interface, renderer render.Renderer, client kai.ProductClient) *Handler {
	return &Handler{
		logger:        logger,
		renderer:      renderer,
		configService: productconfiguration.NewProductConfigService(logger),
		productClient: client,
	}
}
