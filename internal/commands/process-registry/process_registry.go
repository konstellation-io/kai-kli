package processregistry

import (
	registry "github.com/konstellation-io/kli/api/process-registry"
	"github.com/konstellation-io/kli/authserver"
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

func NewHandler(logger logging.Interface, renderer render.Renderer,
	registryCli registry.ProcessRegistryInterface) *Handler {
	return &Handler{
		logger:                logger,
		renderer:              renderer,
		processRegistryClient: registryCli,
		authentication:        auth.NewAuthentication(logger, authserver.NewDefaultAuthServer(logger)),
		configService:         configuration.NewKaiConfigService(logger),
	}
}
