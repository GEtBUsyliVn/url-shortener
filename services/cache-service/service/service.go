package service

import (
	"context"
	"errors"

	mem "github.com/GEtBUsyliVn/url-shortener/services/cache-service/repository/memory"
	"github.com/GEtBUsyliVn/url-shortener/services/cache-service/repository/redis"
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
	v = s.memory.Get(key)

	if v == "" {
		s.log.Info("no cache value in memory, try to get from redis", zap.String("key", key))
		v, err = s.redis.Get(ctx, key)
		if v == "" {
			s.log.Info("no cache value in redis")
		}
	}
	s.log.Info("cache data", zap.String("value", v))
	return v, err
}

func (s *CacheService) Set(ctx context.Context, key, val string) error {
	if err := s.redis.Set(ctx, key, val); err != nil {
		s.log.Info("failed to set cache in redis")
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
