package main

import (
	"context"
	"errors"
	logBasic "log"
	"os"
	"sync"

	"github.com/GEtBUsyliVn/url-shortener/services/url-service/config"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/grpc"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/repository"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/repository/postgres"
	srvc "github.com/GEtBUsyliVn/url-shortener/services/url-service/service"
	invalidatorW "github.com/GEtBUsyliVn/url-shortener/services/url-service/worker"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	defer db.Close()
	repo := postgres.NewRepository(logger, db)
	service := srvc.NewService(logger, repo)
	grpcService := grpc.NewGrpcService(logger, service)
	err = grpcService.Init(cfg.GRPC)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	worker := invalidatorW.NewWorker(repo, logger)
	go worker.Work(ctx, cfg.Worker.Interval, wg)
	logger.Info("worker started")

	logger.Info("server started")

	wg.Wait()
}
