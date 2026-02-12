package repository

import (
	"context"

	"github.com/GEtBUsyliVn/url-shortener/services/url-service/repository/entity"
)

type Storage interface {
	Create(ctx context.Context, url *entity.Url) error
	Get(ctx context.Context, shortCode string) (*entity.Url, error)
	Delete(ctx context.Context, shortCode string) error
	Exists(ctx context.Context, code string) (bool, error)
	UpdateExpired(ctx context.Context) (int, error)
	List(ctx context.Context) ([]*entity.Url, error)
}
