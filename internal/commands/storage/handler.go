package storage

import (
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type Handler struct {
	logger        logging.Interface
	configService *configuration.KaiConfigService
}

func NewHandler(
	logger logging.Interface,
) *Handler {
	return &Handler{
		logger:        logger,
		configService: configuration.NewKaiConfigService(logger),
	}
}
