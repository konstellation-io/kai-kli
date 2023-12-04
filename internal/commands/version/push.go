package version

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/konstellation-io/kli/internal/services/configuration"
)

type PushVersionOpts struct {
	Server      string
	KrtFilePath string
	Version     string
	Description string
}

func (h *Handler) PushVersion(opts *PushVersionOpts) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(opts.Server)
	if err != nil {
		return err
	}

	_, product := h.splitFilePath(opts.KrtFilePath)

	err = h.updateProductVersion(opts)
	if err != nil {
		return err
	}

	versionTag, err := h.versionClient.Push(server, product, opts.KrtFilePath)
	if err != nil {
		return fmt.Errorf("pushing krt.yaml file: %w", err)
	}

	h.logger.Success(fmt.Sprintf("Version with tag %q of product %q successfully created!", versionTag, product))

	return nil
}

func (h *Handler) splitFilePath(krtFilePath string) (dir, filename string) {
	dir, product := filepath.Split(krtFilePath)

	return dir, strings.Split(product, ".")[0]
}

func (h *Handler) updateProductVersion(opts *PushVersionOpts) error {
	if opts.Version == "" {
		return nil
	}

	dir, product := h.splitFilePath(opts.KrtFilePath)

	productConfig, err := h.productService.GetConfiguration(product, dir)
	if err != nil {
		return err
	}

	err = productConfig.UpdateProductVersion(opts.Version, opts.Description)
	if err != nil {
		return err
	}

	err = h.productService.WriteConfiguration(productConfig, product, dir)
	if err != nil {
		return err
	}

	return nil
}
