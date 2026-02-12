package service

import (
	"context"
	"errors"
	"time"

	analyticsC "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc"
	modelA "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/model"
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/model"
	cacheC "github.com/GEtBUsyliVn/url-shortener/services/cache-service/pkg/api/grpc"
	shortenerC "github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/grpc"
	model2 "github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/model"
	"go.uber.org/zap"
)

type GatewayService struct {
	cache     *cacheC.Client
	analytics *analyticsC.Client
	shortener *shortenerC.Client
	log       *zap.Logger
}

func NewGatewayService(cache *cacheC.Client, analytics *analyticsC.Client,
	shortener *shortenerC.Client, log *zap.Logger) *GatewayService {
	return &GatewayService{
		cache:     cache,
		analytics: analytics,
		shortener: shortener,
		log:       log,
	}
}

func (g *GatewayService) CreateShortUrl(ctx context.Context, req *model.ShortCode) (string, error) {
	request := &model2.CreateUrlRequest{
		Url:      req.URL,
		ExpireAt: time.Now().Add(time.Duration(req.ExpiredDays) * 24 * time.Hour),
		UserId:   req.UserId,
	}
	code, err := g.shortener.CreateUrl(ctx, request)
	if err != nil {
		return "", err
	}
	err = g.cache.Set(ctx, req.URL, code)
	if err != nil {
		g.log.Error("failed to set cache", zap.Error(err))
	}
	return code, err
}

func (g *GatewayService) GetOriginalUrl(ctx context.Context, shortCode string, e *modelA.ClickRequest) (string, error) {
	url, err := g.cache.Get(ctx, shortCode)
	if err == nil {
		return url, nil
	}
	g.log.Error("cache miss", zap.String("short_code", shortCode), zap.Error(err))

	url, err = g.shortener.GetOriginalUrl(ctx, shortCode)
	if err != nil {
		if errors.Is(err, shortenerC.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	err = g.cache.Set(ctx, shortCode, url)
	if err != nil {
		g.log.Error("failed to set cache", zap.Error(err))
	}

	err = g.analytics.ClickEvent(ctx, e)
	if err != nil {
		g.log.Error("failed to record click event", zap.Error(err))
	}

	return url, nil

}

func (g *GatewayService) GetAnalytics(ctx context.Context, shortCode string) (*model.Stats, error) {
	analytics, err := g.analytics.GetStatistics(ctx, &modelA.StatsRequest{ShortCode: shortCode})
	if err != nil {
		if errors.Is(err, analyticsC.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resp := &model.Stats{}
	resp.Bind(analytics)
	return resp, err
}
