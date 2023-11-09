package storage

import (
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Handler struct {
	logger        logging.Interface
	renderer      render.Renderer
	configService *configuration.KaiConfigService
}

func NewHandler(
	logger logging.Interface,
	renderer render.Renderer,
) *Handler {
	return &Handler{
		logger:        logger,
		renderer:      renderer,
		configService: configuration.NewKaiConfigService(logger),
	}
}
