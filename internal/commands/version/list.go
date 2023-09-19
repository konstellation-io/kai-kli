package version

type ListVersionsOpts struct {
	ServerName string
	ProductID  string
}

func (c *Handler) ListVersions(opts *ListVersionsOpts) error {
	kaiConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	registeredProcesses, err := c.versionClient.List(srv, opts.ProductID)
	if err != nil {
		return err
	}

	c.renderer.RenderVersions(registeredProcesses)

	return nil
}
