package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Worker struct {
	Interval time.Duration `envconfig:"URL_WORKER_INTERVAL" default:"1m"`
}

func (w *Worker) Prepare(prefix string) error {
	return envconfig.Process(prefix, w)
}
