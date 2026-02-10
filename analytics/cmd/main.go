package main

import (
	"context"
	"fmt"

	"github.com/GEtBUsyliVn/url-shortener/analytics/config"
	"github.com/GEtBUsyliVn/url-shortener/analytics/grpc"
	"github.com/GEtBUsyliVn/url-shortener/analytics/repository"
	"github.com/GEtBUsyliVn/url-shortener/analytics/repository/postgres"
	service2 "github.com/GEtBUsyliVn/url-shortener/analytics/service"
	"github.com/GEtBUsyliVn/url-shortener/analytics/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	router := gin.Default()
	cfg := config.InitConfig()
	fmt.Println(cfg)
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	db, err := repository.NewDB(cfg, logger)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx, cancel := context.WithCancel(context.Background())
	repo := postgres.NewRepository(logger, db)
	service := service2.NewStatsService(repo, logger)
	//aggregator := worker.NewClicksAggregator(repo, logger)
	collector := worker.NewClicksCollector(service, logger)
	collector.Start(ctx)
	grpcServer := grpc.NewGrpcService(logger, service, collector)
	defer cancel()
	//go aggregator.Work(ctx, *logger, time.Second*5)
	grpcServer.Init(cfg.GRPC.Addr)
	router.POST("/url", func(c *gin.Context) {})
	fmt.Println("server started")
	router.Run()
}
