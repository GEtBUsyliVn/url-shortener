package config

import "github.com/kelseyhightower/envconfig"

type GRPC struct {
	Addr string `envconfig:"ANALYTICS_GRPC_ADDR"`
}

func (c *GRPC) Prepare(prefix string) error {
	return envconfig.Process(prefix, c)
}
