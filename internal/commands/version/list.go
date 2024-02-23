package version

import "github.com/konstellation-io/kli/internal/entity"

type ListVersionsOpts struct {
	ServerName   string
	ProductID    string
	StatusFilter string
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

	// TODO - technical debt, this should be done in the server side
	// filter is done here instead of in the client to avoid adding complexity to the client
	// and make the rendering of empty results after filtering more easy to manage
	if opts.StatusFilter != "" {
		registeredProcesses = filterVersionsByStatus(registeredProcesses, opts.StatusFilter)
	}

	h.renderer.RenderVersions(opts.ProductID, registeredProcesses)

	return nil
}

func filterVersionsByStatus(versions []*entity.Version, status string) []*entity.Version {
	var filteredVersions []*entity.Version

	for _, v := range versions {
		if v.Status == status {
			filteredVersions = append(filteredVersions, v)
		}
	}

	return filteredVersions
}
