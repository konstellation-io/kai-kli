package process

import (
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/krt/pkg/krt"

	"github.com/konstellation-io/kli/internal/logging"
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

	NetworkSourcePort *int
	NetworkTargetPort *int
	NetworkProtocol   krt.NetworkingProtocol
	// end of not implemented yet
	Image         string
	Replicas      *int
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

func (w *Handler) getNetwork(opts *ProcessOpts) *krt.ProcessNetworking {
	if opts.NetworkSourcePort != nil && opts.NetworkTargetPort != nil {
		net := &krt.ProcessNetworking{}
		net.TargetPort = *opts.NetworkSourcePort
		net.DestinationPort = *opts.NetworkTargetPort
		net.Protocol = krt.NetworkingProtocolHTTP

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
