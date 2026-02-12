package config

import "github.com/kelseyhightower/envconfig"

type GRPC struct {
	UrlAddr       string `envconfig:"GRPC_URL_ADDR"`
	CacheAddr     string `envconfig:"GRPC_CACHE_ADDR"`
	AnalyticsAddr string `envconfig:"GRPC_ANALYTICS_ADDR"`
}

func (c *GRPC) Prepare(prefix string) error {
	return envconfig.Process(prefix, c)
}
