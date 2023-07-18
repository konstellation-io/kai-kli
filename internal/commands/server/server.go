package server

import (
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/authentication"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Handler struct {
	logger         logging.Interface
	renderer       render.Renderer
	authentication *authentication.AuthenticationService
	configService  *configuration.KaiConfigService
}

func NewServerHandler(logger logging.Interface, renderer render.Renderer) *Handler {
	return &Handler{
		logger:         logger,
		renderer:       renderer,
		authentication: authentication.NewAuthentication(logger),
		configService:  configuration.NewKaiConfigService(logger),
	}
}

type Server struct {
	Name string
	URL  string
}

type RemoteServerInfo struct {
	IsKaiServer bool
}
