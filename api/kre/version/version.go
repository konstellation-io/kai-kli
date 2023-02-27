package version

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kli/api/graphql"
)

// VersionInterface method to interact with Versions.
type VersionInterface interface { //nolint: golint
	List(runtime string) (List, error)
	Start(runtime, versionName, comment string) error
	Stop(runtime, versionName, comment string) error
	Publish(runtime, versionName, comment string) error
	Unpublish(runtime, versionName, comment string) error
	GetConfig(runtime, versionName string) (*Config, error)
	UpdateConfig(runtime, versionName string, configVars []ConfigVariableInput) (bool, error)
	Create(runtime, krtFile string) (string, error)
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
