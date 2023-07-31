package process

import (
	"github.com/konstellation-io/krt/pkg/krt"
)

type AddProcessOpts struct {
	ServerName  string
	ProductID   string
	WorkflowID  string
	ProcessID   string
	ProcessType krt.ProcessType
	// not implemented yet
	ObjectStoreName   string
	ObjectStoreScope  krt.ObjectStoreScope
	NetworkSourcePort int
	NetworkTargetPort int
	NetworkProtocol   krt.NetworkingProtocol
	// end of not implemented yet
	Image         string
	Replicas      int
	GPU           bool
	Subscriptions []string
}

func (w *Handler) AddProcess(opts *AddProcessOpts) error {
	productConfig, err := w.configService.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	// TODO create object store if needed

	// TODO create networking if needed

	err = productConfig.AddProcess(
		opts.WorkflowID,
		&krt.Process{
			Name:          opts.ProcessID,
			Type:          opts.ProcessType,
			Image:         opts.Image,
			Replicas:      &opts.Replicas,
			GPU:           &opts.GPU,
			Config:        map[string]string{},
			ObjectStore:   nil,
			Secrets:       map[string]string{},
			Subscriptions: opts.Subscriptions,
			Networking:    nil,
		})

	if err != nil {
		return err
	}

	err = w.configService.WriteConfiguration(productConfig, opts.ProductID)
	if err != nil {
		return err
	}

	_, wf, err := productConfig.GetWorkflow(opts.WorkflowID)
	if err != nil {
		return err
	}

	w.renderer.RenderProcesses(wf.Processes)

	return nil
}
