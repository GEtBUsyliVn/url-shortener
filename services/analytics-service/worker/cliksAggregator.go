package worker

import (
	"context"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository"
	"go.uber.org/zap"
)

type ClicksAggregator struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewClicksAggregator(repo repository.Repository, logger *zap.Logger) *ClicksAggregator {
	return &ClicksAggregator{
		repo: repo,
		log:  logger,
	}
}

func (c *ClicksAggregator) Work(ctx context.Context, logger zap.Logger, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.log.Info("got tick from ticker")
			c.Aggregate(ctx)
		case <-ctx.Done():
			logger.Info("stopping clicks aggregator")
			return
		}
	}
}

func (c *ClicksAggregator) Aggregate(ctx context.Context) {
	c.log.Info("starting clicks aggregation")

	// Get all short codes that have clicks to aggregate
	shortCodes, err := c.repo.GetUniqClicks(ctx)
	if err != nil {
		c.log.Error("failed to get short codes with clicks", zap.Error(err))
		return
	}

	for _, shortCode := range shortCodes {
		err := c.repo.CreateStats(ctx, shortCode.ShortCode)
		if err != nil {
			c.log.Error("failed to create stats for short code", zap.String("short_code", shortCode.ShortCode), zap.Error(err))
			continue
		}
		c.log.Info("successfully aggregated stats for short code", zap.String("short_code", shortCode.ShortCode))
	}

	c.log.Info("completed clicks aggregation")
}
