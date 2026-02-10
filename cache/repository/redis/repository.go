package redis

import (
	"context"
	"errors"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/cache/repository"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	db      *redis.Client
	expTime time.Duration
}

func NewRedisStorage(redis *redis.Client, expTime time.Duration) *RedisRepository {
	return &RedisRepository{
		db:      redis,
		expTime: expTime,
	}
}

func (r *RedisRepository) GetClient() *redis.Client {
	return r.db
}

func (r *RedisRepository) Set(ctx context.Context, key string, value string) error {
	if err := r.db.Set(ctx, key, value, r.expTime).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.db.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", repository.ErrNotFount
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisRepository) Del(ctx context.Context, key string) error {
	if err := r.db.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
