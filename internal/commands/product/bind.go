package product

import (
	"errors"

	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/pkg/krtutil"
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

	var version *entity.Version
	version, err = h.versionClient.Get(server, opts.ProductID, opts.VersionTag)
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

	err = h.configService.WriteConfiguration(&productconfiguration.KaiProductConfiguration{
		Krt: krtutil.MapVersionToKrt(version),
	}, opts.ProductID, opts.LocalPath)

	if err != nil {
		return err
	}

	h.logger.Success("Product successfully bound!")

	return nil
}
