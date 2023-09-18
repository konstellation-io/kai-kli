package version

import (
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Handler struct {
	logger        logging.Interface
	renderer      render.Renderer
	versionClient kai.VersionClient
	configService *configuration.KaiConfigService
}

func NewHandler(
	logger logging.Interface,
	renderer render.Renderer,
	client kai.VersionClient,
) *Handler {
	return &Handler{
		logger:        logger,
		renderer:      renderer,
		versionClient: client,
		configService: configuration.NewKaiConfigService(logger),
	}
}
