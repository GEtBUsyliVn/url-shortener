package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/GEtBUsyliVn/url-shortener/analytics/model"
	"github.com/GEtBUsyliVn/url-shortener/analytics/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/analytics/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsGrpcService struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	service    *service.BasicService
	proto.UnimplementedAnalyticsServiceServer
	listener net.Listener
}

func NewGrpcService(log *zap.Logger, service *service.BasicService) *AnalyticsGrpcService {
	return &AnalyticsGrpcService{
		log:        log.Named("grpc analytics"),
		service:    service,
		grpcServer: grpc.NewServer(),
	}
}

func (s *AnalyticsGrpcService) Init(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("unable to listen port %w", err)
	}
	s.listener = listener

	proto.RegisterAnalyticsServiceServer(s.grpcServer, s)

	go func() {
		s.log.Debug("grpc serve", zap.String("address", addr))

		if servErr := s.grpcServer.Serve(s.listener); servErr != nil {
			s.log.Panic("unable to serve", zap.Error(servErr))
		}
	}()

	s.log.Info("server started", zap.String("addr", addr))

	return nil
}

func (s *AnalyticsGrpcService) RecordClick(ctx context.Context, req *proto.ClickEvent) (*emptypb.Empty, error) {
	click := &model.Click{}
	click.BindProto(req)
	if err := s.service.CreateClick(ctx, click); err != nil {
		s.log.Error("failed to create click", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create click: %v", err)
	}
	return nil, nil
}
