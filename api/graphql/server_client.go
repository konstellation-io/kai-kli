package graphql

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/konstellation-io/kli/api/kai/config"

	"github.com/konstellation-io/graphql"
)

// GqlManager struct to implement access to GraphQL endpoints with gql client.
type GqlManager struct {
	appVersion string
	cfg        *ClientConfig
	server     *config.ServerConfig
	client     *graphql.Client
	httpClient *http.Client
}

func (g *GqlManager) setupClient(args ...graphql.ClientOption) error {
	if g.client != nil {
		return nil
	}

	accessToken, err := g.getAccessToken()
	if err != nil {
		return err
	}

	c := NewHTTPClient([]Option{
		AddHeader("User-Agent", "Konstellation KLI"),
		AddHeader("KLI-Version", g.appVersion),
		AddHeader("Cache-Control", "no-cache"),
		AddHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)),
	}...)

	opts := []graphql.ClientOption{graphql.WithHTTPClient(c)}
	opts = append(opts, args...)

	g.client = graphql.NewClient(fmt.Sprintf("%s/graphql", g.server.URL), opts...)
	g.httpClient = c

	if g.cfg.Debug {
		g.client.Log = func(s string) { log.Println(s) }
	}

	return nil
}

// MakeRequest call to GraphQL endpoint.
func (g *GqlManager) MakeRequest(query string, vars map[string]interface{}, respData interface{}) error {
	err := g.setupClient()
	if err != nil {
		return err
	}

	req := graphql.NewRequest(query)

	for k, v := range vars {
		req.Var(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.cfg.DefaultRequestTimeout)
	defer cancel()

	err = g.client.Run(ctx, req, respData)
	if err != nil {
		return fmt.Errorf("graphql error: %w", err)
	}

	return nil
}

// UploadFile uploads a file to KAI server.
func (g *GqlManager) UploadFile(file graphql.File, query string, vars map[string]interface{}, respData interface{}) error {
	err := g.setupClient(graphql.UseMultipartForm())
	if err != nil {
		return err
	}

	req := graphql.NewRequest(query)

	ctx, cancel := context.WithTimeout(context.Background(), g.cfg.DefaultRequestTimeout)
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

// NewGqlManager creates an instance of GqlManager that takes cares of authentication.
func NewGqlManager(cfg *ClientConfig, server *config.ServerConfig, appVersion string) *GqlManager {
	return &GqlManager{
		appVersion,
		cfg,
		server,
		nil,
		nil,
	}
}
