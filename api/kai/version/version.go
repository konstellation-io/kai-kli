package version

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kli/api/graphql"
)

// VersionInterface method to interact with Versions.
type VersionInterface interface { //nolint: revive
	List(product string) (List, error)
	Start(product, versionName, comment string) error
	Stop(product, versionName, comment string) error
	Publish(product, versionName, comment string) error
	Unpublish(product, versionName, comment string) error
	GetConfig(product, versionName string) (*Config, error)
	UpdateConfig(product, versionName string, configVars []ConfigVariableInput) (bool, error)
	Create(product, krtFile string) (string, error)
}

type versionClient struct {
	client *graphql.GqlManager
}

// New creates a new struct to access Versions methods.
func New(gql *graphql.GqlManager) VersionInterface {
	return &versionClient{
		gql,
	}
}
