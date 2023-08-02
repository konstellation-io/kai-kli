package process

import (
	"github.com/konstellation-io/krt/pkg/krt"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

//nolint:revive // the name needs to be ProcessOpts
type ProcessOpts struct {
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

type Handler struct {
	logger        logging.Interface
	renderer      render.Renderer
	configService *productconfiguration.ProductConfigService
}

func NewHandler(logger logging.Interface, renderer render.Renderer) *Handler {
	return &Handler{
		logger:        logger,
		renderer:      renderer,
		configService: productconfiguration.NewProductConfigService(logger),
	}
}

func (w *Handler) getCPULimits(opts *ProcessOpts) *krt.ResourceLimit {
	res := &krt.ResourceLimit{}

	res.Request = _defaultCPURequest
	res.Limit = _defaultCPULimit

	if opts.CPURequest != "" {
		res.Request = opts.CPURequest
		res.Limit = res.Request
	}

	if opts.CPULimit != "" {
		res.Limit = opts.CPULimit
	}

	return res
}

func (w *Handler) getMemoryLimits(opts *ProcessOpts) *krt.ResourceLimit {
	res := &krt.ResourceLimit{}

	res.Request = _defaultMemoryRequest
	res.Limit = _defaultMemoryLimit

	if opts.MemoryRequest != "" {
		res.Request = opts.MemoryRequest
		res.Limit = res.Request
	}

	if opts.MemoryLimit != "" {
		res.Limit = opts.MemoryLimit
	}

	return res
}

func (w *Handler) getNetwork(opts *ProcessOpts) *krt.ProcessNetworking {
	if opts.NetworkSourcePort != -1 && opts.NetworkTargetPort != -1 {
		net := &krt.ProcessNetworking{}
		net.TargetPort = opts.NetworkSourcePort
		net.DestinationPort = opts.NetworkTargetPort
		net.Protocol = krt.NetworkingProtocolTCP

		if opts.NetworkProtocol != "" {
			net.Protocol = opts.NetworkProtocol
		}

		return net
	}

	return nil
}

func (w *Handler) getObjectStore(opts *ProcessOpts) *krt.ProcessObjectStore {
	if opts.ObjectStoreName != "" {
		obj := &krt.ProcessObjectStore{}
		obj.Name = opts.ObjectStoreName
		obj.Scope = krt.ObjectStoreScopeWorkflow

		if opts.ObjectStoreScope != "" {
			obj.Scope = opts.ObjectStoreScope
		}

		return obj
	}

	return nil
}
