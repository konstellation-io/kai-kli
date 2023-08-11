package server

func (c *Handler) SetDefaultServer(server string) error {
	kaiConfig, err := c.configService.GetConfiguration()

	if err != nil {
		return err
	}

	err = kaiConfig.SetDefaultServer(server)
	if err != nil {
		return err
	}

	err = c.configService.WriteConfiguration(kaiConfig)
	if err != nil {
		return err
	}

	c.renderer.RenderServers(kaiConfig.Servers)

	return nil
}
