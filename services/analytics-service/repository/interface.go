package repository

import (
	"context"

	entity2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository/entity"
)

type Repository interface {
	CreateClick(ctx context.Context, click *entity2.Click) error
	CreateStats(ctx context.Context, shortCode string) error
	GetStatistics(ctx context.Context, shortCode string) (*entity2.Statistics, error)
	GetUniqClicks(ctx context.Context) ([]*entity2.Click, error)
	CreateClicksBatch(ctx context.Context, clicks []*entity2.Click) error
}
