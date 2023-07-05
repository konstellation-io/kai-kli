package kre

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/kreclient.go -package=mocks

import (
	"github.com/konstellation-io/kli/api/kre/product"
	"github.com/konstellation-io/kli/api/kre/version"
)

// KreInterface first level methods.
type KreInterface interface { //nolint:revive
	Version() version.VersionInterface
	Product() product.ProductInterface
}
