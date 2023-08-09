package authserver

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (as *AuthServer) buildAuthorizationRequest() string {
	return fmt.Sprintf(
		"%v/realms/%v/protocol/openid-connect/auth?client_id=%v&redirect_uri=%v&response_mode=query&response_type=code&scope=openid",
		as.config.KeycloakConfig.KeycloakURL,
		as.config.KeycloakConfig.Realm,
		as.config.KeycloakConfig.ClientID,
		as.getCallbackURL(),
	)
}

func (as *AuthServer) buildTokenExchangeRequest(code string) (*http.Request, error) {
	tokenURL := fmt.Sprintf("%v/realms/%v/protocol/openid-connect/token",
		as.config.KeycloakConfig.KeycloakURL,
		as.config.KeycloakConfig.Realm)

	body := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"client_id":    {as.config.KeycloakConfig.ClientID},
		"redirect_uri": {as.getCallbackURL()},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}