package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/GEtBUsyliVn/url-shortener/url/model"
	"github.com/GEtBUsyliVn/url-shortener/url/repository"
	"go.uber.org/zap"
)

type Service struct {
	repo repository.Storage
	log  *zap.Logger
}

func NewService(log *zap.Logger, repo repository.Storage) *Service {
	return &Service{
		repo: repo,
		log:  log.Named("service"),
	}
}

func (s *Service) CreateShortURL(ctx context.Context, url *model.Url) (string, error) {
	code, err := GenerateShortCode(func(code string) (bool, error) {
		return s.repo.Exists(ctx, code)
	})

	if err != nil {
		return "", err
	}

	ent := url.Entity()
	ent.ShortCode = code
	err = s.repo.Create(ctx, ent)
	if err != nil {
		return "", fmt.Errorf("failed to create url: %w", err)
	}
	return code, nil
}

func (s *Service) GetShortUrl(ctx context.Context, shortUrl string) (string, error) {
	entity, err := s.repo.Get(ctx, shortUrl)
	if errors.Is(err, repository.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}
	return entity.OriginalUrl, nil
}
