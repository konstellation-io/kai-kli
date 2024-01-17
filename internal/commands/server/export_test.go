//go:build unit

package server

import (
	"github.com/konstellation-io/kli/authserver"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/auth"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func NewTestHandler(logger logging.Interface, renderer render.Renderer,
	authServer authserver.Authenticator) *Handler {
	return &Handler{
		logger:         logger,
		renderer:       renderer,
		authentication: auth.NewAuthentication(logger, authServer),
		configService:  configuration.NewKaiConfigService(logger),
	}
}
