package service

import (
	"context"
	"errors"

	mem "github.com/GEtBUsyliVn/url-shortener/cache/repository/memory"
	"github.com/GEtBUsyliVn/url-shortener/cache/repository/redis"
	"go.uber.org/zap"
)

type CacheService struct {
	redis  *redis.RedisRepository
	memory *mem.MemoryRepository
	log    *zap.Logger
}

func NewCacheService(redis *redis.RedisRepository, memory *mem.MemoryRepository, logger *zap.Logger) *CacheService {
	return &CacheService{
		redis:  redis,
		memory: memory,
		log:    logger.Named("cache service"),
	}
}

func (s *CacheService) Get(ctx context.Context, key string) (string, error) {
	var v string
	var err error
	val := s.memory.Get(key)

	if val == "" {
		v, err = s.redis.Get(ctx, key)
	}

	return v, err
}

func (s *CacheService) Set(ctx context.Context, key, val string) error {
	if err := s.redis.Set(ctx, key, val); err != nil {
		return err
	}
	s.memory.Set(key, val)
	return nil
}

func (s *CacheService) Del(ctx context.Context, key string) error {
	if err := s.redis.Del(ctx, key); err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	s.memory.Del(key)
	return nil
}
