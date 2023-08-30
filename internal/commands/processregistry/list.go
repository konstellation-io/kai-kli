package processregistry

import (
	"fmt"

	"github.com/konstellation-io/krt/pkg/krt"
)

type ListProcessesOpts struct {
	ServerName  string
	ProductID   string
	ProcessType krt.ProcessType
}

func (c *Handler) ListProcesses(opts *ListProcessesOpts) error {
	if opts.ProcessType != "" && !opts.ProcessType.IsValid() {
		return fmt.Errorf("invalid process type: %q", opts.ProcessType)
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
		srv, opts.ProductID, string(opts.ProcessType),
	)
	if err != nil {
		return err
	}

	c.logger.Info(fmt.Sprintf("Registered processes for %q product:\n", opts.ProductID))
	c.renderer.RenderRegisteredProcesses(registeredProcesses)

	return nil
}
