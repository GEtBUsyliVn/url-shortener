package cacheCleaner

import (
	"context"
	"time"

	mem "github.com/GEtBUsyliVn/url-shortener/cache/repository/memory"
	"go.uber.org/zap"
)

type Worker struct {
	repo *mem.MemoryRepository
	log  *zap.Logger
}

func NewWorker(repo *mem.MemoryRepository, log *zap.Logger) *Worker {
	return &Worker{
		repo: repo,
		log:  log,
	}
}

func (w *Worker) Work(duration time.Duration, ctx context.Context) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.ClearCache()
		}
	}
}

func (w *Worker) ClearCache() {
	n := w.repo.DeleteExpired(time.Now())
	if n > 0 {
		w.log.Info("expired cache deleted", zap.Int("count", n))
	}
}
