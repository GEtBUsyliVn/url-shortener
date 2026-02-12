package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/model"
	proto2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc/proto"
	service2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/service"
	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/worker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsGrpcService struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	service    *service2.BasicService
	collector  *worker.ClicksCollector
	proto2.UnimplementedAnalyticsServiceServer
	listener net.Listener
}

func NewGrpcService(log *zap.Logger, service *service2.BasicService, collector *worker.ClicksCollector) *AnalyticsGrpcService {
	return &AnalyticsGrpcService{
		log:        log.Named("grpc analytics"),
		service:    service,
		grpcServer: grpc.NewServer(),
		collector:  collector,
	}
}

func (s *AnalyticsGrpcService) Init(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("unable to listen port %w", err)
	}
	s.listener = listener

	proto2.RegisterAnalyticsServiceServer(s.grpcServer, s)

	go func() {
		s.log.Debug("grpc serve", zap.String("address", addr))

		if servErr := s.grpcServer.Serve(s.listener); servErr != nil {
			s.log.Panic("unable to serve", zap.Error(servErr))
		}
	}()

	s.log.Info("grpc server started", zap.String("addr", addr))

	return nil
}

func (s *AnalyticsGrpcService) RecordClick(ctx context.Context, req *proto2.ClickEvent) (*emptypb.Empty, error) {
	click := &model.Click{}
	click.BindProto(req)
	ok := s.collector.TryEnqueue(ctx, click)
	if !ok {
		s.log.Error("clicks collector is full, dropping click event", zap.Any("click", click))
	}
	return &emptypb.Empty{}, nil
}

func (s *AnalyticsGrpcService) GetStatistics(ctx context.Context, req *proto2.StatsRequest) (*proto2.StatsResponse, error) {
	stats, err := s.service.GetStats(ctx, req.ShortCode)
	if err != nil {
		if errors.Is(err, service2.ErrNotFound) {
			s.log.Error("failed to get stats", zap.Error(err))
			return nil, status.Errorf(codes.NotFound, "unknown short code: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get stats: %v", err)
	}
	return stats.BindProtoResponse(), nil
}

func (s *AnalyticsGrpcService) Shutdown() {
	s.grpcServer.Stop()
	s.log.Info("grpc shut downed")
}
