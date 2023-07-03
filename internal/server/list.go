package server

import "fmt"

func (c *KaiConfigurator) ListServers() error {
	kaiConfiguration, err := c.getConfiguration()
	if err != nil {
		return fmt.Errorf("failed to get configuration: %w", err)
	}

	c.renderer.RenderServers(kaiConfiguration)

	return nil
}
