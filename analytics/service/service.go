package service

import (
	"context"
	"errors"

	"github.com/GEtBUsyliVn/url-shortener/analytics/model"
	"github.com/GEtBUsyliVn/url-shortener/analytics/repository"
	"github.com/GEtBUsyliVn/url-shortener/analytics/repository/entity"
	"go.uber.org/zap"
)

type BasicService struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewStatsService(repo repository.Repository, log *zap.Logger) *BasicService {
	return &BasicService{
		repo: repo,
		log:  log.Named("analytics service"),
	}
}

func (s *BasicService) CreateClick(ctx context.Context, clicks []*model.Click) {
	entities := make([]*entity.Click, 0)
	for _, click := range clicks {
		entities = append(entities, click.Entity())
	}
	if err := s.repo.CreateClicksBatch(ctx, entities); err != nil {
		s.log.Error("batch insert failed", zap.Error(err))
	}
}

func (s *BasicService) GetStats(ctx context.Context, shortCode string) (*model.Statistics, error) {
	stats, err := s.repo.GetStatistics(ctx, shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	stat := &model.Statistics{}
	stat.Bind(stats)

	return stat, nil
}
