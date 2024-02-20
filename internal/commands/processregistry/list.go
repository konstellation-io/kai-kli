package processregistry

import (
	"fmt"

	"github.com/konstellation-io/krt/pkg/krt"
)

type ListProcessesOpts struct {
	ServerName string
	ProductID  string
	Filters    *ListFilters
}

type ListFilters struct {
	ProcessName string
	Version     string
	ProcessType krt.ProcessType
}

func (c *Handler) ListProcesses(opts *ListProcessesOpts) error {
	if opts.Filters.ProcessType != "" && !opts.Filters.ProcessType.IsValid() {
		return fmt.Errorf("invalid process type: %q", opts.Filters.ProcessType)
	}

	kaiConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	registeredProcesses, err := c.processRegistryClient.List(
		srv, opts.ProductID, opts.Filters.ProcessName, opts.Filters.Version, string(opts.Filters.ProcessType),
	)
	if err != nil {
		return err
	}

	c.renderer.RenderRegisteredProcesses(registeredProcesses)

	return nil
}
