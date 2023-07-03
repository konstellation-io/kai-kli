package kre_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/api/kre/config"
	config2 "github.com/konstellation-io/kli/cmd/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func gqlMockServer(t *testing.T, requestVars, mockResponse string) (*httptest.Server, *graphql.GqlManager) {
	t.Helper()

	auth := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		stringBody := string(b)
		if !auth {
			require.Equal(t, stringBody, `{"apiToken":"12345"}`)
			auth = true
			_, err = io.WriteString(w, `{"access_token": "access_12345"}`)
			require.NoError(t, err)
			return
		}

		if requestVars != "" {
			actualBody := map[string]interface{}{}
			err := json.NewDecoder(strings.NewReader(stringBody)).Decode(&actualBody)
			require.NoError(t, err)

			expectedVars := map[string]interface{}{}
			err = json.NewDecoder(strings.NewReader(requestVars)).Decode(&expectedVars)
			require.NoError(t, err)

			require.EqualValues(t, expectedVars, actualBody["variables"])
		}

		if mockResponse == "" {
			mockResponse = "{}"
		}
		_, err = io.WriteString(w, mockResponse)
		require.NoError(t, err)
	}))

	viper.Set("request_timeout", 999999*time.Second)
	srvCfg := &config.ServerConfig{
		Name:     "test",
		URL:      srv.URL,
		APIToken: "12345",
	}

	clientConfig := &graphql.ClientConfig{
		DefaultRequestTimeout: viper.GetDuration("request_timeout"),
		Debug:                 viper.GetBool(config2.DebugKey),
	}
	client := graphql.NewGqlManager(clientConfig, srvCfg, "test")

	return srv, client
}
