package version

import "github.com/konstellation-io/kli/api/version"

type GetVersionOpts struct {
	ServerName string
	ProductID  string
	VersionTag string
}

func (h *Handler) GetVersion(opts *GetVersionOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	registeredProcess, err := h.versionClient.Get(srv, opts.ProductID, opts.VersionTag)
	if err != nil {
		return err
	}

	h.renderer.RenderVersions([]*version.Version{registeredProcess})

	return nil
}
