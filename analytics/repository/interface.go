package repository

import (
	"context"

	"github.com/GEtBUsyliVn/url-shortener/analytics/repository/entity"
)

type Repository interface {
	CreateClick(ctx context.Context, click *entity.Click) error
	CreateStats(ctx context.Context, shortCode string) error
	GetStatistics(ctx context.Context, shortCode string) (*entity.Statistics, error)
	GetUniqClicks(ctx context.Context) ([]*entity.Click, error)
	CreateClicksBatch(ctx context.Context, clicks []*entity.Click) error
}
