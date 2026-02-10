package main

import (
	"github.com/GEtBUsyliVn/url-shortener/url/config"
	"github.com/GEtBUsyliVn/url-shortener/url/grpc"
	grpcClient "github.com/GEtBUsyliVn/url-shortener/url/pkg/api/grpc"
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
	defer db.Close()
	repo := postgres.NewRepository(logger, db)
	service := srvc.NewService(logger, repo)
	grpcService := grpc.NewGrpcService(logger, service)
	err = grpcService.Init(cfg.Grpc)
	if err != nil {
		panic(err)
	}
	_ = grpcClient.NewGrpcClient(cfg.Grpc.Addr, false, logger)
	router := gin.Default()

	//router.POST("/url", func(c *gin.Context) {
	//	req := &model.CreateUrlRequest{}
	//	if err := c.ShouldBindJSON(req); err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//		return
	//	}
	//	url, err := client.CreateUrl(c, req)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		return
	//	}
	//	c.JSON(http.StatusOK, gin.H{"code": url})
	//})
	//
	//router.GET("/url/:code", func(c *gin.Context) {
	//	code := c.Param("code")
	//	if code == "" {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
	//		return
	//	}
	//	url, err := service.GetShortUrl(c, code)
	//	if errors.Is(err, repository.ErrNotFound) {
	//		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
	//		return
	//	}
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		return
	//	}
	//	c.Header("Location", url)
	//	c.Redirect(http.StatusMovedPermanently, "https://github.com")
	//})
	router.GET("/health", func(c *gin.Context) {})

	router.Run(":8080")

}
