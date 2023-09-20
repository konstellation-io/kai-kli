package version

import "fmt"

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

	tag, err := h.versionClient.Start(
		srv, opts.ProductID, opts.VersionTag, opts.Comment,
	)
	if err != nil {
		return err
	}

	h.logger.Success(
		fmt.Sprintf("Version %q is starting.", tag),
	)

	return nil
}