package version

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/internal/config"
)

// VersionInterface method to interact with Versions.
type VersionInterface interface { //nolint: golint
	List() (List, error)
	Start(versionName, comment string) error
	Stop(versionName, comment string) error
	Publish(versionName, comment string) error
	Unpublish(versionName, comment string) error
	GetConfig(versionName string) (*Config, error)
	UpdateConfig(versionName string, configVars []ConfigVariableInput) (bool, error)
	Create(krtFile string) (string, error)
}

type versionClient struct {
	cfg    *config.Config
	client *graphql.GqlManager
}

// New creates a new struct to access Versions methods.
func New(cfg *config.Config, gql *graphql.GqlManager) VersionInterface {
	return &versionClient{
		cfg,
		gql,
	}
}
