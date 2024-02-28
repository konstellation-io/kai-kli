package product

import (
	"errors"
	"strings"

	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/pkg/krtutil"
)

var (
	ErrProductExistsNoVersions = errors.New("error product exists but has no versions")
)

type BindProductOpts struct {
	Server     string
	ProductID  string
	VersionTag *string
	LocalPath  string
	Force      bool
}

func (h *Handler) Bind(opts *BindProductOpts) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(opts.Server)
	if err != nil {
		return err
	}

	productConfig, err := h.configService.GetConfiguration(opts.ProductID, opts.LocalPath)
	if err != nil && !errors.Is(err, productconfiguration.ErrProductConfigNotFound) {
		return err
	}

	if productConfig != nil && !opts.Force {
		return ErrProductConfigurationAlreadyExists
	}

	version, err := h.getVersion(server, opts)
	if err != nil {
		return err
	}

	err = h.configService.WriteConfiguration(&productconfiguration.KaiProductConfiguration{
		Krt: krtutil.MapVersionToKrt(version),
	}, opts.ProductID, opts.LocalPath)

	if err != nil {
		return err
	}

	h.renderer.RenderProductBinded(opts.ProductID)

	return nil
}

func (h *Handler) getVersion(server *configuration.Server, opts *BindProductOpts) (*entity.Version, error) {
	version, err := h.versionClient.Get(server, opts.ProductID, opts.VersionTag)
	if err != nil {
		if strings.Contains(err.Error(), ErrProductExistsNoVersions.Error()) {
			h.logger.Warn("Product exists but has no versions, using an empty version")
			version = &entity.Version{}
		} else {
			return nil, err
		}
	}

	return version, nil
}
