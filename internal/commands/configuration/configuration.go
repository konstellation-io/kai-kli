package productconfiguration

import (
	"errors"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

var (
	ErrScopeNotValid = errors.New("scope not valid, valid scopes are: version, workflow, process")
)

const (
	_versionScope  = "version"
	_workflowScope = "workflow"
	_processScope  = "process"
)

type Handler struct {
	logger        logging.Interface
	renderer      render.Renderer
	productConfig *productconfiguration.ProductConfigService
}

func NewHandler(logger logging.Interface, renderer render.Renderer) *Handler {
	return &Handler{
		logger:        logger,
		renderer:      renderer,
		productConfig: productconfiguration.NewProductConfigService(logger),
	}
}
