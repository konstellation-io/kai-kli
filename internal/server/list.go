package server

import "fmt"

func (c *Handler) ListServers() error {
	kaiConfiguration, err := c.configHandler.GetConfiguration()
	if err != nil {
		return fmt.Errorf("failed to get configuration: %w", err)
	}

	c.renderer.RenderServers(kaiConfiguration.Servers)

	return nil
}
