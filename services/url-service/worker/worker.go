package worker

import (
	"context"
	"sync"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/url-service/repository"
	"go.uber.org/zap"
)

type Worker struct {
	repo repository.Storage
	log  *zap.Logger
}

func NewWorker(repo repository.Storage, logger *zap.Logger) *Worker {
	return &Worker{
		repo: repo,
		log:  logger,
	}
}

func (w *Worker) Work(ctx context.Context, duration time.Duration, wg *sync.WaitGroup) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			count, err := w.repo.UpdateExpired(ctx)
			if err != nil {
				w.log.Error("failed to update expired urls", zap.Error(err))
			}

			w.log.Info("updated expired urls", zap.Int("count", count))
		}
	}
}
