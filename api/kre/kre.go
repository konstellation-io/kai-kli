package kre

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kli/api/kre/version"
)

// KreInterface first level methods.
type KreInterface interface { //nolint: golint
	Version() version.VersionInterface
}
