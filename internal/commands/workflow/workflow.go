package workflow

import (
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

type KaiWorkflow struct {
	logger        logging.Interface
	renderer      render.Renderer
	configService *productconfiguration.ProductConfigService
}

func NewHandler(logger logging.Interface, renderer render.Renderer) *KaiWorkflow {
	return &KaiWorkflow{
		logger:        logger,
		renderer:      renderer,
		configService: productconfiguration.NewProductConfigService(logger),
	}
}
