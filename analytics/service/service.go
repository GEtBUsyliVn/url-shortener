package service

import (
	"context"
	"errors"

	"github.com/GEtBUsyliVn/url-shortener/analytics/model"
	"github.com/GEtBUsyliVn/url-shortener/analytics/repository"
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

func (s *BasicService) CreateClick(ctx context.Context, click *model.Click) error {
	err := s.repo.CreateClick(ctx, click.Entity())
	if err != nil {
		s.log.Error("failed to create click", zap.Error(err))
		return err
	}
	return nil
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
