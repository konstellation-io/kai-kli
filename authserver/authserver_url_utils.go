package authserver

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (as *AuthServer) buildAuthorizationRequest(config KeycloakConfig) string {
	return fmt.Sprintf(
		"%v/realms/%v/protocol/openid-connect/auth?client_id=%v&redirect_uri=%v&response_mode=query&response_type=code&scope=openid",
		config.KeycloakURL,
		config.Realm,
		config.ClientID,
		as.getCallbackURL(),
	)
}

func (as *AuthServer) buildTokenExchangeRequest(code string, config KeycloakConfig) (*http.Request, error) {
	tokenURL := fmt.Sprintf("%v/realms/%v/protocol/openid-connect/token",
		config.KeycloakURL,
		config.Realm)

	body := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"client_id":    {config.ClientID},
		"redirect_uri": {as.getCallbackURL()},
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, tokenURL, strings.NewReader(body.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, err
}
