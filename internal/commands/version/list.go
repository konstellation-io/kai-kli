package version

type ListVersionsOpts struct {
	ServerName string
	ProductID  string
}

func (h *Handler) ListVersions(opts *ListVersionsOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	registeredProcesses, err := h.versionClient.List(srv, opts.ProductID)
	if err != nil {
		return err
	}

	h.renderer.RenderVersions(opts.ProductID, registeredProcesses)

	return nil
}
