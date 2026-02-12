package main

import (
	"context"
	"errors"
	logBasic "log"
	"os"
	"sync"

	"github.com/GEtBUsyliVn/url-shortener/services/cache-service/cacheCleaner"
	"github.com/GEtBUsyliVn/url-shortener/services/cache-service/config"
	"github.com/GEtBUsyliVn/url-shortener/services/cache-service/grpc"
	"github.com/GEtBUsyliVn/url-shortener/services/cache-service/repository/memory"
	redisS "github.com/GEtBUsyliVn/url-shortener/services/cache-service/repository/redis"
	"github.com/GEtBUsyliVn/url-shortener/services/cache-service/service"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Overload(); err != nil {
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			logBasic.Fatal(err)
		}
	}

	cfg := &config.Config{}
	if err := cfg.Prepare(config.AppName); err != nil {
		logBasic.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password, // no password set
		DB:       cfg.Redis.Db,       // use default DB
	})
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	redisStorage := redisS.NewRedisStorage(redisClient, cfg.Redis.CacheTtl)
	memStorage := memory.NewMemoryStorage(cfg.Memory.CacheTtl)

	zapCfg := zap.NewProductionConfig()
	zapCfg.DisableStacktrace = true
	logger, err := zapCfg.Build()

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	cacheService := service.NewCacheService(redisStorage, memStorage, logger)
	grpcService := grpc.NewGrpcService(logger, cacheService)
	worker := cacheCleaner.NewWorker(memStorage, logger)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go worker.Work(cfg.Worker.Interval, ctx, wg)

	err = grpcService.Init(cfg.GRPC.Addr)
	if err != nil {
		panic(err)
	}
	wg.Wait()
	defer grpcService.Shutdown()
	logger.Info("analytics server started")
	defer redisStorage.GetClient().Close()

}
