package server

import (
	"github.com/konstellation-io/kli/internal/configuration"
)

func (c *Handler) SetDefaultServer(server string) error {
	configHandler := configuration.NewKaiConfigHandler(c.logger)
	kaiConfig, err := configHandler.GetConfiguration()

	if err != nil {
		return err
	}

	err = kaiConfig.SetDefaultServer(server)
	if err != nil {
		return err
	}

	err = configHandler.WriteConfiguration(kaiConfig)
	if err != nil {
		return err
	}

	c.renderer.RenderServers(kaiConfig.Servers)

	return nil
}
