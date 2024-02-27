package product

import "github.com/konstellation-io/kli/internal/services/configuration"

func (h *Handler) RemoveUser(serverName, product, userEmail string) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(serverName)
	if err != nil {
		return err
	}

	err = h.productClient.RemoveUserFromProduct(server, product, userEmail)
	if err != nil {
		return err
	}

	h.renderer.RenderRemoveUserFromProduct(product, userEmail)

	return nil
}
