package worker

import (
	"context"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/analytics/model"
	"github.com/GEtBUsyliVn/url-shortener/analytics/service"
	"go.uber.org/zap"
)

type ClicksCollector struct {
	ch      chan *model.Click
	service *service.BasicService
	log     *zap.Logger
}

func NewClicksCollector(svc *service.BasicService, logger *zap.Logger) *ClicksCollector {
	return &ClicksCollector{
		ch:      make(chan *model.Click, 1000),
		service: svc,
		log:     logger,
	}
}

// не блокируем хендлер, уважаем ctx (отменили запрос — не пушим)
func (c *ClicksCollector) TryEnqueue(ctx context.Context, click *model.Click) bool {
	select {
	case <-ctx.Done():
		return false
	case c.ch <- click:
		c.log.Info("wrote click to channel")
		return true
	default:
		return false
	}
}

// запускать 1 раз
func (c *ClicksCollector) Start(ctx context.Context) {
	c.log.Info("starting clicks collector")
	const (
		maxBatch      = 100
		flushInterval = 1 * time.Second
	)

	go func() {
		batch := make([]*model.Click, 0, maxBatch)
		ticker := time.NewTicker(flushInterval)
		defer ticker.Stop()

		flush := func() {
			if len(batch) == 0 {
				return
			}
			// ВАЖНО: не ctx из запроса; используем ctx воркера
			c.service.CreateClick(ctx, batch)
			batch = batch[:0]
		}

		for {
			select {
			case <-ctx.Done():
				flush()
				return

			case ev, ok := <-c.ch:
				if !ok {
					flush()
					return
				}
				if ev == nil {
					continue
				}

				batch = append(batch, ev)
				if len(batch) >= maxBatch {
					c.log.Info("clicks batch full, flushing", zap.Int("batch_size", len(batch)))
					flush()
				}

			case <-ticker.C:
				c.log.Info("batch flush interval reached, flushing", zap.Int("batch_size", len(batch)))
				flush()
			}
		}
	}()
}
