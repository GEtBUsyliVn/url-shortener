package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Worker struct {
	Interval time.Duration `envconfig:"WORKER_CLEANUP_INTERVAL"`
}

func (w *Worker) Prepare(prefix string) error {
	return envconfig.Process(prefix, w)
}
