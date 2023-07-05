package kai

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/kaiclient.go -package=mocks -mock_names Client=MockKaiClient

import (
	"github.com/konstellation-io/kli/api/kai/product"
	"github.com/konstellation-io/kli/api/kai/version"
)

// Client first level methods.
type Client interface {
	Version() version.VersionInterface
	Product() product.ProductInterface
}
