package main

import (
	"context"
	"fmt"

	"github.com/GEtBUsyliVn/url-shortener/cache/cacheCleaner"
	"github.com/GEtBUsyliVn/url-shortener/cache/config"
	"github.com/GEtBUsyliVn/url-shortener/cache/grpc"
	"github.com/GEtBUsyliVn/url-shortener/cache/repository/memory"
	redisS "github.com/GEtBUsyliVn/url-shortener/cache/repository/redis"
	"github.com/GEtBUsyliVn/url-shortener/cache/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	cfg := config.InitConfig()
	ctx, cancel := context.WithCancel(context.Background())
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password, // no password set
		DB:       cfg.Redis.DB,       // use default DB
	})
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	redisStorage := redisS.NewRedisStorage(redisClient, cfg.Redis.CacheTTl)
	memStorage := memory.NewMemoryStorage(cfg.MemoryStorage.CacheTTl)
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	cacheService := service.NewCacheService(redisStorage, memStorage, logger)
	grpcService := grpc.NewGrpcService(logger, cacheService)
	worker := cacheCleaner.NewWorker(memStorage, logger)
	go worker.Work(cfg.Worker.Interval, ctx)
	err = grpcService.Init(cfg.Grpc.Addr)
	if err != nil {
		panic(err)
	}
	defer grpcService.Shutdown()
	fmt.Println("server started")
	defer redisStorage.GetClient().Close()

}
