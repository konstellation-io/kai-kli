package kai_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/konstellation-io/kli/api/kai/config"

	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/graphql"
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

	cfg := &config.Config{
		DefaultRequestTimeout: 999999 * time.Second,
	}
	srvCfg := &config.ServerConfig{
		Name:     "test",
		URL:      srv.URL,
		APIToken: "12345",
	}

	clientConfig := &graphql.ClientConfig{
		DefaultRequestTimeout: cfg.DefaultRequestTimeout,
		Debug:                 cfg.Debug,
	}
	client := graphql.NewGqlManager(clientConfig, srvCfg, "test")

	return srv, client
}
