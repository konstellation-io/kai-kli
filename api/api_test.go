package api_test

import (
	"testing"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKreClient(t *testing.T) {
	d := testhelpers.SetupConfigDir(t)
	defer testhelpers.CleanConfigDir(t, d)

	cfg, err := config.NewConfig("")
	assert.NoError(t, err)

	clientCfg := &graphql.ClientConfig{
		DefaultRequestTimeout: viper.GetDuration("request_timeout"),
		Debug:                 cfg.Debug,
	}

	srv := config.ServerConfig{
		Name:     "test",
		URL:      "http://test",
		APIToken: "12345",
	}
	err = cfg.AddServer(srv)
	require.NoError(t, err)

	k := api.NewKreClient(clientCfg, &srv, "test-version")

	require.NotEmpty(t, k.Version())
}
