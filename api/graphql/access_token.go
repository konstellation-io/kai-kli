package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// accessTokenResponse represents a response from the sign-in API.
type accessTokenResponse struct {
	Token string `json:"access_token"`
}

// getAccessToken call to sign-in endpoint and get an access_token to use in later API calls.
func (g *GqlManager) getAccessToken() (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/auth/token/signin", g.server.URL)

	postData := bytes.NewBuffer([]byte(fmt.Sprintf(`{"apiToken":%q}`, g.server.APIToken)))

	ctx, cancel := context.WithTimeout(context.Background(), g.cfg.DefaultRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, postData)
	if err != nil {
		return "", fmt.Errorf("error creating request call: %s", err) //nolint:goerr113
	}

	req.Header.Set("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling access token URL: %s", err) //nolint:goerr113
	}
	defer r.Body.Close()

	if r.Body == nil {
		return "", ErrResponseEmpty
	}

	var t accessTokenResponse

	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		return "", fmt.Errorf("error decoding access token response: %s", err) //nolint:goerr113
	}

	return t.Token, nil
}
