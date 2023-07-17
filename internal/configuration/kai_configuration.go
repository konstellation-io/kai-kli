package configuration

import (
	"time"
)

type KaiConfiguration struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	Realm     string `yaml:"realm"`
	ClientID  string `yaml:"clientId"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Token     *Token `yaml:"token"`
	IsDefault bool   `yaml:"default"`
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
	return t.Date.Add(time.Duration(t.ExpiresIn) * time.Second).After(time.Now().UTC())
}
