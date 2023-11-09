package storage

import (
	"fmt"

	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/pkg/osutil"
)

func (h *Handler) Login(serverName string) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(serverName)
	if err != nil {
		return err
	}

	if err := osutil.OpenBrowser(server.StorageURL); err != nil {
		h.logger.Warn(fmt.Sprintf("Unable to open browser, open the following URL: %v",
			server.StorageURL))
		return nil
	}

	h.logger.Success("Storage login page's opening in the browser!")
	return nil
}
