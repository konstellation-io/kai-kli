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
	CPURequest    string
	CPULimit      string
	MemoryRequest string
	MemoryLimit   string

	ObjectStoreName  string
	ObjectStoreScope krt.ObjectStoreScope

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

	obj := w.getObjectStore(opts)
	network := w.getNetwork(opts)
	CPULimits := w.getCPULimits(opts)
	memoryLimits := w.getMemoryLimits(opts)

	limits := &krt.ProcessResourceLimits{}

	if CPULimits != nil || memoryLimits != nil {
		limits.CPU = CPULimits
		limits.Memory = memoryLimits
	}

	err = productConfig.AddProcess(
		opts.WorkflowID,
		&krt.Process{
			Name:           opts.ProcessID,
			Type:           opts.ProcessType,
			Image:          opts.Image,
			Replicas:       &opts.Replicas,
			GPU:            &opts.GPU,
			Config:         map[string]string{},
			ObjectStore:    obj,
			Secrets:        map[string]string{},
			Subscriptions:  opts.Subscriptions,
			ResourceLimits: limits,
			Networking:     network,
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

func (w *Handler) getCPULimits(opts *AddProcessOpts) *krt.ResourceLimit {
	res := &krt.ResourceLimit{}

	if opts.CPURequest != "" {
		res.Request = opts.CPURequest
		res.Limit = res.Request

		if opts.CPULimit != "" {
			res.Limit = opts.CPULimit
		}
	}

	return res
}

func (w *Handler) getMemoryLimits(opts *AddProcessOpts) *krt.ResourceLimit {
	res := &krt.ResourceLimit{}

	if opts.MemoryRequest != "" {
		res.Request = opts.MemoryRequest
		res.Limit = res.Request

		if opts.MemoryLimit != "" {
			res.Limit = opts.MemoryLimit
		}
	}

	return res
}

func (w *Handler) getNetwork(opts *AddProcessOpts) *krt.ProcessNetworking {
	net := &krt.ProcessNetworking{}

	if opts.NetworkSourcePort != 0 && opts.NetworkTargetPort != 0 {
		net.TargetPort = opts.NetworkSourcePort
		net.DestinationPort = opts.NetworkTargetPort
		net.Protocol = krt.NetworkingProtocolTCP

		if opts.NetworkProtocol != "" {
			net.Protocol = opts.NetworkProtocol
		}
	}

	return net
}

func (w *Handler) getObjectStore(opts *AddProcessOpts) *krt.ProcessObjectStore {
	obj := &krt.ProcessObjectStore{}
	if opts.ObjectStoreName != "" {
		obj.Name = opts.ObjectStoreName
		obj.Scope = krt.ObjectStoreScopeWorkflow

		if opts.ObjectStoreScope != "" {
			obj.Scope = opts.ObjectStoreScope
		}
	}

	return obj
}
