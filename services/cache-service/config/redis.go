package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Redis struct {
	Addr     string        `envconfig:"REDIS_ADDR"`
	Password string        `envconfig:"REDIS_PASSWORD"`
	Db       int           `envconfig:"REDIS_DB"`
	CacheTtl time.Duration `envconfig:"REDIS_CACHE_TTL"`
}

func (c *Redis) Prepare(prefix string) error {
	return envconfig.Process(prefix, c)
}
