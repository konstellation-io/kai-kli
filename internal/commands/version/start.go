package version

import (
	"github.com/konstellation-io/kli/internal/entity"
)

type StartOpts struct {
	ServerName string
	ProductID  string
	VersionTag string
	Comment    string
}

func (h *Handler) Start(opts *StartOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	vers, err := h.versionClient.Get(srv, opts.ProductID, &opts.VersionTag)
	if err != nil {
		return err
	}

	if vers.Status == entity.VersionStatusCritical {
		h.renderer.RenderCallout(vers)
	}

	tag, err := h.versionClient.Start(
		srv, opts.ProductID, opts.VersionTag, opts.Comment,
	)
	if err != nil {
		return err
	}

	h.renderer.RenderStartVersion(opts.ProductID, tag)

	return nil
}
