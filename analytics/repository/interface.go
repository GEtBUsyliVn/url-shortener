package repository

import (
	"context"

	"github.com/GEtBUsyliVn/url-shortener/analytics/repository/entity"
)

type Repository interface {
	CreateClick(ctx context.Context, click *entity.Click) error
	CreateStats(ctx context.Context, shortCode string) error
	GetStatistics(ctx context.Context, shortCode string) (*entity.Statistics, error)
}
