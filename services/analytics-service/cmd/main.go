package main

import (
	"context"
	"errors"
	logBasic "log"
	"os"
	"sync"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/config"
	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/grpc"
	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository"
	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository/postgres"
	service2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/service"
	worker2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/worker"
	"github.com/joho/godotenv"
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

	zapCfg := zap.NewProductionConfig()
	zapCfg.DisableStacktrace = true
	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	db, err := repository.NewDB(cfg, logger)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	repo := postgres.NewRepository(logger, db)
	service := service2.NewStatsService(repo, logger)
	aggregator := worker2.NewClicksAggregator(repo, logger)
	collector := worker2.NewClicksCollector(service, logger)
	wg.Add(1)
	collector.Start(ctx, wg)
	grpcServer := grpc.NewGrpcService(logger, service, collector)
	defer cancel()
	go aggregator.Work(ctx, *logger, cfg.Worker.Interval)
	grpcServer.Init(cfg.GRPC.Addr)
	defer grpcServer.Shutdown()
	logger.Info("analytics server started")
	wg.Wait()
}
