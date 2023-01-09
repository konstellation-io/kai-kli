package kre_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/graphql"
	"github.com/konstellation-io/kli/internal/config"
)

func gqlMockServer(t *testing.T, requestVars, mockResponse string) (*httptest.Server, *config.Config, *graphql.GqlManager) {
	t.Helper()

	auth := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body) //nolint:gocritic
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

			require.EqualValues(t, actualBody["variables"], expectedVars)
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

	client := graphql.NewGqlManager(cfg, srvCfg, "test")

	return srv, cfg, client
}
