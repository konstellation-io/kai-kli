package server

func (c *Handler) RemoveServer(server string) error {
	configService, err := c.configService.GetConfiguration()
	if err != nil {
		return err
	}

	err = configService.DeleteServer(server)
	if err != nil {
		return err
	}

	err = c.configService.WriteConfiguration(configService)
	if err != nil {
		return err
	}

	c.renderer.RenderServers(configService.Servers)

	return nil
}
