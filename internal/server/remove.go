package server

import "github.com/konstellation-io/kli/internal/configuration"

func (c *Handler) RemoveServer(server string) error {
	config, err := c.configHandler.GetConfiguration()
	if err != nil {
		return err
	}

	if !config.CheckServerExists(server) {
		return configuration.ErrServerNotFound
	}

	err = config.DeleteServer(server)
	if err != nil {
		return err
	}

	err = c.configHandler.WriteConfiguration(config)
	if err != nil {
		return err
	}

	c.renderer.RenderServers(config.Servers)

	return nil
}
