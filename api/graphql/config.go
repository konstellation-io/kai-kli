package graphql

import "time"

type ClientConfig struct {
	DefaultRequestTimeout time.Duration
	Debug                 bool
}
