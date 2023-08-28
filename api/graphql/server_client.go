package graphql

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/konstellation-io/graphql"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

// GqlManager struct to implement access to GraphQL endpoints with gql client.
type GqlManager struct {
	appVersion string
	client     *graphql.Client
	httpClient *http.Client
}

// NewGqlManager creates an instance of GqlManager that takes cares of authentication.
func NewGqlManager() *GqlManager {
	return &GqlManager{
		viper.GetString(config.BuildVersionKey),
		nil,
		nil,
	}
}

// MakeRequest call to GraphQL endpoint.
func (g *GqlManager) MakeRequest(server *configuration.Server, query string, vars map[string]interface{},
	respData interface{}) error {
	err := g.setupClient(server)
	if err != nil {
		return err
	}

	req := graphql.NewRequest(query)

	for k, v := range vars {
		req.Var(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.RequestTimeoutKey))
	defer cancel()

	err = g.client.Run(ctx, req, respData)
	if err != nil {
		return fmt.Errorf("graphql error: %w", err)
	}

	return nil
}

// UploadFile uploads a file to KAI server.
func (g *GqlManager) UploadFile(server *configuration.Server, file graphql.File, query string,
	vars map[string]interface{}, respData interface{}) error {
	err := g.setupClient(server, graphql.UseMultipartForm())
	if err != nil {
		return err
	}

	req := graphql.NewRequest(query)

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.RequestTimeoutKey))
	defer cancel()

	req.File(file.Field, file.Name, file.R)

	for k, v := range vars {
		req.Var(k, v)
	}

	err = g.client.Run(ctx, req, respData)
	if err != nil {
		return fmt.Errorf("graphql error: %w", err)
	}

	return nil
}

func (g *GqlManager) setupClient(server *configuration.Server, args ...graphql.ClientOption) error {
	if g.client != nil {
		return nil
	}

	clientOpts := []Option{
		AddHeader("User-Agent", "Konstellation KLI"),
		AddHeader("KLI-Version", g.appVersion),
		AddHeader("Cache-Control", "no-cache"),
	}

	if server.IsLoggedIn() {
		clientOpts = append(clientOpts,
			AddHeader("Authorization", fmt.Sprintf("Bearer %s", server.Token.AccessToken)))
	}

	c := NewHTTPClient(clientOpts...)

	opts := []graphql.ClientOption{graphql.WithHTTPClient(c)}
	opts = append(opts, args...)

	g.client = graphql.NewClient(fmt.Sprintf("%s/graphql", server.URL), opts...)
	g.httpClient = c

	if viper.GetBool(config.DebugKey) {
		g.client.Log = func(s string) { log.Println(s) }
	}

	return nil
}
