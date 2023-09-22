package version

import "fmt"

type PublishOpts struct {
	ServerName string
	ProductID  string
	VersionTag string
	Comment    string
}

func (h *Handler) Publish(opts *PublishOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	publishedTriggers, err := h.versionClient.Publish(
		srv, opts.ProductID, opts.VersionTag, opts.Comment,
	)
	if err != nil {
		return err
	}

	if len(publishedTriggers) == 0 {
		h.logger.Info("No triggers found to publish.")
		return nil
	}

	h.logger.Success(
		fmt.Sprintf("%q - %q correctly published.", opts.ProductID, opts.VersionTag),
	)

	h.renderer.RenderTriggers(publishedTriggers)

	return nil
}
