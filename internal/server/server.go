package server

import (
	"github.com/konstellation-io/kli/internal/configuration"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

type Handler struct {
	logger        logging.Interface
	renderer      render.Renderer
	configHandler *configuration.KaiConfigHandler
}

func NewServerHandler(logger logging.Interface, renderer render.Renderer) *Handler {
	return &Handler{
		logger:        logger,
		renderer:      renderer,
		configHandler: configuration.NewKaiConfigHandler(logger),
	}
}

type Server struct {
	Name string
	URL  string
}

type RemoteServerInfo struct {
	IsKaiServer bool
}
