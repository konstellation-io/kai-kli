package product

import "github.com/konstellation-io/kli/internal/services/configuration"

func (h *Handler) ListProducts(serverName string) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(serverName)
	if err != nil {
		return err
	}

	products, err := h.productClient.GetProducts(server)
	if err != nil {
		return err
	}

	h.renderer.RenderProducts(products)

	return nil
}
