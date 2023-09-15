package version

import (
	"github.com/konstellation-io/kli/api/graphql"
)

type Client struct {
	gqlClient *graphql.GqlManager
}

func NewClient(gqlClient *graphql.GqlManager) *Client {
	return &Client{
		gqlClient: gqlClient}
}
