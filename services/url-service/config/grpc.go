package config

import "github.com/kelseyhightower/envconfig"

type GRPC struct {
	Addr string `envconfig:"URL_SERVICE_GRPC_ADDR"`
}

func (g *GRPC) Prepare(prefix string) error {
	return envconfig.Process(prefix, g)
}
