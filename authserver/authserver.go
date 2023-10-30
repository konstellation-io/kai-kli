package authserver

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/phayes/freeport"
)

const (
	_defaultRandomPort   = 0
	_defaultCallbackPath = "sso-callback"
)

//go:embed static/success.html
var successPage string

//go:generate mockgen -source=${GOFILE} -destination=../mocks/auth_server.go -package=mocks AuthServer
type AuthServerInterface interface { //nolint:revive
	StartServer(config KeycloakConfig) (*AuthResponse, error)
}

type AuthServer struct {
	config   EmbeddedServerConfig
	closeApp sync.WaitGroup
	logger   logging.Interface
	response *AuthResponse
}

func NewDefaultAuthServer(logger logging.Interface) *AuthServer {
	return NewAuthServer(logger, EmbeddedServerConfig{
		Port:         _defaultRandomPort,
		CallbackPath: _defaultCallbackPath,
	})
}

func NewAuthServer(logger logging.Interface, config EmbeddedServerConfig) *AuthServer {
	if config.Port <= 0 {
		port, err := freeport.GetFreePort()
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting free port: %s", err))

			config.Port = 3000
		}

		config.Port = uint32(port)
	}

	return &AuthServer{
		logger:   logger,
		config:   config,
		closeApp: sync.WaitGroup{},
	}
}

type KeycloakConfig struct {
	KeycloakURL string
	Realm       string
	ClientID    string
}

type EmbeddedServerConfig struct {
	Port         uint32
	CallbackPath string
}

type AuthResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
}

func (as *AuthServer) openBrowser(url string) error {
	var browserCommand *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		browserCommand = exec.Command("xdg-open", url)
	case "windows":
		browserCommand = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		browserCommand = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}

	err := browserCommand.Run()

	return err
}

func (as *AuthServer) getCallbackURL() string {
	return fmt.Sprintf("http://localhost:%v/%v",
		as.config.Port,
		as.config.CallbackPath)
}

func (as *AuthServer) StartServer(config KeycloakConfig) (*AuthResponse, error) {
	serverAddress := fmt.Sprintf("localhost:%v", as.config.Port)

	as.logger.Debug(fmt.Sprintf("Booting up the server at: %s", serverAddress))

	as.closeApp.Add(1)

	http.HandleFunc(fmt.Sprintf("/%s", as.config.CallbackPath),
		func(w http.ResponseWriter, r *http.Request) {
			as.logger.Debug(fmt.Sprintf("Callback received: %v", r.URL))

			code := r.URL.Query().Get("code")
			if code == "" {
				as.logger.Info("Code not found in the callback URL")
				as.closeApp.Done()
				return
			}

			tokenResponse, err := as.tokenExchangeRequest(code, config)
			if err != nil {
				as.logger.Error(fmt.Sprintf("Error exchanging token: %v", err))
				as.closeApp.Done()
				return
			}

			as.response = tokenResponse

			_, err = fmt.Fprint(w, successPage)
			if err != nil {
				as.logger.Error(fmt.Sprintf("Error writing response: %v", err))
			}

			as.closeApp.Done()
		})

	go func() {
		server := http.Server{
			Addr:              serverAddress,
			ReadHeaderTimeout: 5 * time.Minute,
		}

		if err := server.ListenAndServe(); err != nil {
			as.logger.Error(fmt.Sprintf("Unable to start server: %v\n", err))
			as.closeApp.Done()
		}
	}()

	err := as.openBrowser(as.buildAuthorizationRequest(config))
	if err != nil {
		as.logger.Warn(fmt.Sprintf("Unable to open browser, open the following URL: %v",
			as.buildAuthorizationRequest(config)))
	}

	as.closeApp.Wait()

	if as.response != nil {
		return as.response, nil
	}

	return nil, fmt.Errorf("unable to get access token")
}

func (as *AuthServer) tokenExchangeRequest(code string, config KeycloakConfig) (*AuthResponse, error) {
	request, err := as.buildTokenExchangeRequest(code, config)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange code for token: %v", err)
	}

	var resp *http.Response

	var body []byte

	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange code for token: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		content, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		switch content {
		case "application/json":
			var tokenResponse AuthResponse

			err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
			if err != nil {
				return nil, fmt.Errorf("error decoding token response: %v", err)
			}

			return &tokenResponse, nil
		default:
			return nil, fmt.Errorf("unexpected content type: %v", body)
		}
	}

	return nil, fmt.Errorf("invalid Status code (%v)", resp.StatusCode)
}
