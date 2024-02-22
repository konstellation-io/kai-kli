package product

import (
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type ListProductsOpts struct {
	ProductName string
}

func (h *Handler) ListProducts(serverName string, opts *ListProductsOpts) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(serverName)
	if err != nil {
		return err
	}

	if opts == nil {
		opts = &ListProductsOpts{}
	}

	products, err := h.productClient.GetProducts(server, opts.ProductName)
	if err != nil {
		return err
	}

	h.renderer.RenderProducts(products)

	return nil
}

func filterProductsByName(products []kai.Product, name string) []kai.Product {
	var filteredProducts []kai.Product

	for _, p := range products {
		if p.Name == name {
			filteredProducts = append(filteredProducts, p)
		}
	}

	return filteredProducts
}
