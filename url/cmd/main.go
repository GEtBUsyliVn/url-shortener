package main

import (
	"net/http"

	"github.com/GEtBUsyliVn/url-shortener/url/config"
	"github.com/GEtBUsyliVn/url-shortener/url/grpc"
	grpcClient "github.com/GEtBUsyliVn/url-shortener/url/pkg/api/grpc"
	"github.com/GEtBUsyliVn/url-shortener/url/pkg/api/model"
	"github.com/GEtBUsyliVn/url-shortener/url/repository"
	"github.com/GEtBUsyliVn/url-shortener/url/repository/postgres"
	srvc "github.com/GEtBUsyliVn/url-shortener/url/service"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	cfg := config.InitConfig()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	db, err := repository.NewDB(cfg, logger)
	if err != nil {
		panic(err)
	}
	repo := postgres.NewRepository(logger, db)
	service := srvc.NewService(logger, repo)
	grpcService := grpc.NewGrpcService(logger, service)
	err = grpcService.Init(cfg.Grpc)
	if err != nil {
		panic(err)
	}
	client := grpcClient.NewGrpcClient(cfg.Grpc.Port, false, logger)
	router := gin.Default()
	router.POST("/url", func(c *gin.Context) {
		req := &model.CreateUrlRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		url, err := client.CreateUrl(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"url": url})
	})

	router.Run(":8000")

}
