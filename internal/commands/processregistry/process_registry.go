package processregistry

import (
	"github.com/konstellation-io/kli/api/processregistry"
	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/auth"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Handler struct {
	logger                logging.Interface
	renderer              render.Renderer
	processRegistryClient processregistry.APIClient
	authentication        *auth.AuthenticationService
	configService         *configuration.KaiConfigService
}

func NewHandler(logger logging.Interface, renderer render.Renderer,
	processRegistryCli processregistry.APIClient) *Handler {
	return &Handler{
		logger:                logger,
		renderer:              renderer,
		processRegistryClient: processRegistryCli,
		authentication:        auth.NewAuthentication(logger, authserver.NewDefaultAuthServer(logger)),
		configService:         configuration.NewKaiConfigService(logger),
	}
}
