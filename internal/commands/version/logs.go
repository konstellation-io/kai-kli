package version

import "github.com/konstellation-io/kli/internal/entity"

type GetLogsOpts struct {
	ServerName    string
	LogFilters    entity.LogFilters
	OutFormat     entity.LogOutFormat
	ShowAllLabels bool
}

func (h *Handler) GetLogs(opts *GetLogsOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	versionLogs, err := h.versionClient.Logs(srv, opts.LogFilters)
	if err != nil {
		return err
	}

	h.renderer.RenderLogs(versionLogs, opts.OutFormat, opts.ShowAllLabels)

	return nil
}
