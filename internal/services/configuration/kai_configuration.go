package configuration

import (
	"time"
)

type KaiConfiguration struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name       string `yaml:"name"`
	Host       string `yaml:"host"`
	Protocol   string `yaml:"protocol"`
	APIURL     string `yaml:"apiUrl"`
	AuthURL    string `yaml:"authUrl"`
	StorageURL string `yaml:"storageUrl"`
	Realm      string `yaml:"realm"`
	ClientID   string `yaml:"clientId"`
	Token      *Token `yaml:"token"`
	IsDefault  bool   `yaml:"default"`
}

func (s *Server) IsLoggedIn() bool {
	return s.Token != nil && s.Token.AccessToken != ""
}

type Token struct {
	Date             time.Time `yaml:"date"`
	AccessToken      string    `yaml:"access_token"`
	ExpiresIn        int       `yaml:"expires_in"`
	RefreshExpiresIn int       `yaml:"refresh_expires_in"`
	RefreshToken     string    `yaml:"refresh_token"`
	TokenType        string    `yaml:"token_type"`
}

func (t *Token) IsValid() bool {
	return t.AccessToken != "" && t.Date.Add(time.Duration(t.ExpiresIn)*time.Second).After(time.Now().UTC())
}
