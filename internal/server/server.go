package server

import (
	auth "github.com/konstellation-io/kli/internal/authentication"
	"github.com/konstellation-io/kli/internal/configuration"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

type Handler struct {
	logger         logging.Interface
	renderer       render.Renderer
	authentication *auth.Authentication
	configService  *configuration.KaiConfigService
}

func NewServerHandler(logger logging.Interface, renderer render.Renderer) *Handler {
	return &Handler{
		logger:         logger,
		renderer:       renderer,
		authentication: auth.NewAuthentication(logger),
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
