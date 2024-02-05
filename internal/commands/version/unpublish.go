package version

type UnpublishOpts struct {
	ServerName string
	ProductID  string
	VersionTag string
	Comment    string
}

func (h *Handler) Unpublish(opts *UnpublishOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	tag, err := h.versionClient.Unpublish(
		srv, opts.ProductID, opts.VersionTag, opts.Comment,
	)
	if err != nil {
		return err
	}

	h.renderer.RenderUnpublishVersion(opts.ProductID, tag)

	return nil
}
