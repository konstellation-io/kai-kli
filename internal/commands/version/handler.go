package version

import (
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

type Handler struct {
	logger        logging.Interface
	renderer      render.Renderer
	versionClient kai.VersionClient
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
	}
}
