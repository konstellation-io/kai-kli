package version

type StopOpts struct {
	ServerName string
	ProductID  string
	VersionTag string
	Comment    string
}

func (h *Handler) Stop(opts *StopOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	tag, err := h.versionClient.Stop(
		srv, opts.ProductID, opts.VersionTag, opts.Comment,
	)
	if err != nil {
		return err
	}

	h.renderer.RenderStopVersion(tag, opts.ProductID)

	return nil
}
