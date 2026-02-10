package main

import (
	cacheGrpc "github.com/GEtBUsyliVn/url-shortener/cache/pkg/api/grpc"
	"github.com/GEtBUsyliVn/url-shortener/gateway/config"
	urlGrpc "github.com/GEtBUsyliVn/url-shortener/url/pkg/api/grpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.InitConfig()
	c := gin.Default()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	_ = cacheGrpc.NewGrpcClient(cfg.Grpc.CacheAddr, false, logger)
	_ = urlGrpc.NewGrpcClient(cfg.Grpc.ShortenerAddr, false, logger)
	c.Run()

}
