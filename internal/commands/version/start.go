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

	tag, status, err := h.versionClient.Start(
		srv, opts.ProductID, opts.VersionTag, opts.Comment,
	)
	if err != nil {
		return err
	}

	h.logger.Success(
		fmt.Sprintf("Success, product version %q status is: %q.", tag, status),
	)

	return nil
}
