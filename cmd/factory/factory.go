package factory

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/kre"
	"github.com/konstellation-io/kli/cmd/iostreams"
	"github.com/konstellation-io/kli/cmd/krttools"
	"github.com/konstellation-io/kli/internal/config"
	"github.com/konstellation-io/kli/internal/logger"
	"github.com/konstellation-io/kli/pkg/errors"
)

// CmdFactory interface to define all methods needed during commands creation.
type CmdFactory interface {
	IOStreams() *iostreams.IOStreams
	Config() *config.Config
	Logger() logger.Interface
	KreClient(string) (kre.KreInterface, error)
	Krt() krttools.KrtTooler
}

// Factory contains all data needed during commands execution.
type Factory struct {
	appVersion string
	ioStreams  *iostreams.IOStreams
	cfg        *config.Config
	logger     logger.Interface
}

// NewCmdFactory returns a new Factory object used during commands execution.
func NewCmdFactory(appVersion string) *Factory {
	io := iostreams.System()

	return &Factory{
		appVersion: appVersion,
		ioStreams:  io,
		cfg:        config.NewConfig(),
		logger:     logger.NewDefaultLogger(),
	}
}

// IOStreams access to IOStreams object.
func (f *Factory) IOStreams() *iostreams.IOStreams {
	return f.ioStreams
}

// Config access to Config object.
func (f *Factory) Config() *config.Config {
	return f.cfg
}

// Logger access to Logger object.
func (f *Factory) Logger() logger.Interface {
	return f.logger
}

// KreClient generates a new ServerClient specific for the given server name.
func (f *Factory) KreClient(serverName string) (kre.KreInterface, error) {
	server := f.cfg.GetByServerName(serverName)
	if server == nil {
		return nil, errors.ErrUnknownServerName
	}

	return api.NewKreClient(f.cfg, server, f.appVersion), nil
}

// Krt tools to build and validate.
func (f *Factory) Krt() krttools.KrtTooler {
	return krttools.NewKrtTools(f.logger)
}
