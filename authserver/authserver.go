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

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/phayes/freeport"
)

//go:embed static/success.html
var successPage string

type AuthServer struct {
	config   Config
	closeApp sync.WaitGroup
	logger   logging.Interface
	response *AuthResponse
}

func NewAuthServer(logger logging.Interface, config Config) *AuthServer {
	if config.EmbeddedServerConfig.Port <= 0 {
		port, err := freeport.GetFreePort()
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting free port: %s", err))
			config.EmbeddedServerConfig.Port = 3000
		}
		config.EmbeddedServerConfig.Port = uint32(port)
	}

	return &AuthServer{
		logger:   logger,
		config:   config,
		closeApp: sync.WaitGroup{},
	}
}

type Config struct {
	KeycloakConfig       KeycloakConfig
	EmbeddedServerConfig EmbeddedServerConfig
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
		as.config.EmbeddedServerConfig.Port,
		as.config.EmbeddedServerConfig.CallbackPath)
}

func (as *AuthServer) StartServer() (*AuthResponse, error) {

	serverAddress := fmt.Sprintf("localhost:%v", as.config.EmbeddedServerConfig.Port)

	as.logger.Info(fmt.Sprintf("Booting up the server at: %s", serverAddress))

	as.closeApp.Add(1)

	http.HandleFunc(fmt.Sprintf("/%s", as.config.EmbeddedServerConfig.CallbackPath),
		func(w http.ResponseWriter, r *http.Request) {
			as.logger.Info(fmt.Sprintf("Callback received: %v", r.URL))

			code := r.URL.Query().Get("code")
			if code != "" {
				request, err := as.buildTokenExchangeRequest(code)
				if err == nil {
					var resp *http.Response
					var body []byte
					resp, err = http.DefaultClient.Do(request)
					if err == nil {
						defer resp.Body.Close()
						if resp.StatusCode == http.StatusOK {
							content, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
							switch content {
							case "application/json":
								var tokenResponse AuthResponse
								err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
								if err != nil {
									as.logger.Error(fmt.Sprintf("Error decoding token response: %v", err))
								}

								as.response = &tokenResponse
							default:
								as.logger.Warn(fmt.Sprintf("Unexpected content type: %v", body))
							}
						} else {
							as.logger.Error(fmt.Sprintf("invalid Status code (%v)", resp.StatusCode))
						}
						_, err := fmt.Fprintf(w, successPage)
						if err != nil {
							as.logger.Error(fmt.Sprintf("Error writing response: %v", err))
						}
						as.closeApp.Done()
						return
					}
					as.logger.Error(fmt.Sprintf("Unable to exchange code for token: %v", err))
				}
				as.logger.Error(fmt.Sprintf("Unable to exchange code for token: %v", err))
			}

			as.logger.Info("Code not found in the callback URL")
		})

	go func() {
		if err := http.ListenAndServe(serverAddress, nil); err != nil {
			as.logger.Error(fmt.Sprintf("Unable to start server: %v\n", err))
			as.closeApp.Done()
		}
	}()

	err := as.openBrowser(as.buildAuthorizationRequest())
	if err != nil {
		return nil, fmt.Errorf("could not open the browser for url %v", as.buildAuthorizationRequest())
	}

	as.closeApp.Wait()

	if as.response != nil {
		return as.response, nil
	}

	return nil, fmt.Errorf("unable to get access token")
}
