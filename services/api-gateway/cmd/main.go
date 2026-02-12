package main

import (
	"errors"
	logBasic "log"
	"os"

	analyticsGrpc "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc"
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/config"
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/rest/handler"
	route "github.com/GEtBUsyliVn/url-shortener/services/api-gateway/router"
	gatewayService "github.com/GEtBUsyliVn/url-shortener/services/api-gateway/service"
	cacheGrpc "github.com/GEtBUsyliVn/url-shortener/services/cache-service/pkg/api/grpc"
	urlGrpc "github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/grpc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	c := gin.Default()
	validate := validator.New(validator.WithRequiredStructEnabled())

	zapCfg := zap.NewProductionConfig()
	zapCfg.DisableStacktrace = true

	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	cache := cacheGrpc.NewGrpcClient(cfg.GRPC.CacheAddr, false, logger)
	shortener := urlGrpc.NewGrpcClient(cfg.GRPC.UrlAddr, false, logger)
	analytics := analyticsGrpc.NewGrpcClient(cfg.GRPC.AnalyticsAddr, false, logger)
	service := gatewayService.NewGatewayService(cache, analytics, shortener, logger)
	urlHandler := handler.NewHandler(logger, service, validate)
	router := route.NewRouter(c, urlHandler)
	router.RegisterRoutes()
	logger.Info("server started")
	c.Run()

}
