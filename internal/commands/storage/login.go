package storage

import (
	"fmt"

	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/pkg/osutil"
)

func (h *Handler) OpenConsole(serverName string) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(serverName)
	if err != nil {
		return err
	}

	if err := osutil.OpenBrowser(server.StorageEndpoint); err != nil {
		h.logger.Warn(fmt.Sprintf("Unable to open browser, open the following URL: %v",
			server.StorageEndpoint))
		return nil
	}

	h.logger.Success("Storage console page's opening in the browser!")

	return nil
}
