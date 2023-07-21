package processregistry

import (
	registry "github.com/konstellation-io/kli/api/process-registry"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/auth"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Handler struct {
	logger                logging.Interface
	renderer              render.Renderer
	processRegistryClient registry.ProcessRegistryInterface
	authentication        *auth.AuthenticationService
	configService         *configuration.KaiConfigService
}

func NewProcessRegistryHandler(logger logging.Interface, renderer render.Renderer,
	registryCli registry.ProcessRegistryInterface, authService *auth.AuthenticationService,
	configService *configuration.KaiConfigService) *Handler {
	return &Handler{
		logger:                logger,
		renderer:              renderer,
		processRegistryClient: registryCli,
		authentication:        authService,
		configService:         configService,
	}
}
