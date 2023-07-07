package kai

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/kaiclient.go -package=mocks

import (
	"github.com/konstellation-io/kli/api/kai/product"
	"github.com/konstellation-io/kli/api/kai/version"
)

// KaiInterface first level methods.
type KaiInterface interface { //nolint: golint
	Version() version.VersionInterface
	Product() product.ProductInterface
}
