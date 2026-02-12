package service

import (
	"context"
	"errors"

	model2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/model"
	repository2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository"
	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository/entity"
	"go.uber.org/zap"
)

type BasicService struct {
	repo repository2.Repository
	log  *zap.Logger
}

func NewStatsService(repo repository2.Repository, log *zap.Logger) *BasicService {
	return &BasicService{
		repo: repo,
		log:  log.Named("analytics service"),
	}
}

func (s *BasicService) CreateClick(ctx context.Context, clicks []*model2.Click) {
	entities := make([]*entity.Click, 0)
	for _, click := range clicks {
		entities = append(entities, click.Entity())
	}
	if err := s.repo.CreateClicksBatch(ctx, entities); err != nil {
		s.log.Error("batch insert failed", zap.Error(err))
	}

}

func (s *BasicService) GetStats(ctx context.Context, shortCode string) (*model2.Statistics, error) {
	stats, err := s.repo.GetStatistics(ctx, shortCode)
	if err != nil {
		if errors.Is(err, repository2.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	stat := &model2.Statistics{}
	stat.Bind(stats)

	return stat, nil
}
