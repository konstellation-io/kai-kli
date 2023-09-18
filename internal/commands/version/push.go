package version

import (
	"fmt"
	"path"
	"strings"

	"github.com/konstellation-io/kli/internal/services/configuration"
)

type PushVersionOpts struct {
	Server      string
	KrtFilePath string
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

	product := h.getProductFromFilePath(opts.KrtFilePath)

	versionTag, err := h.versionClient.Push(server, product, opts.KrtFilePath)
	if err != nil {
		return fmt.Errorf("pushing krt.yaml file: %w", err)
	}

	h.logger.Success(fmt.Sprintf("Version with tag %q of product %q successfully created!", versionTag, product))

	return nil
}

func (h *Handler) getProductFromFilePath(krtFilePath string) string {
	return strings.Split(path.Base(krtFilePath), ".")[0]
}
