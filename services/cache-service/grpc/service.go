package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	proto2 "github.com/GEtBUsyliVn/url-shortener/services/cache-service/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/services/cache-service/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	log          *zap.Logger
	grpcServer   *grpc.Server
	cacheService *service.CacheService
	proto2.UnimplementedCacheServiceServer
	listener net.Listener
}

func NewGrpcService(log *zap.Logger, cacheService *service.CacheService) *Service {
	return &Service{
		log:          log.Named("grpc cache"),
		cacheService: cacheService,
		grpcServer:   grpc.NewServer(),
	}
}

func (s *Service) Init(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("unable to listen port %w", err)
	}
	s.listener = listener

	proto2.RegisterCacheServiceServer(s.grpcServer, s)

	go func() {
		s.log.Debug("grpc serve", zap.String("address", addr))

		if servErr := s.grpcServer.Serve(s.listener); servErr != nil {
			s.log.Panic("unable to serve", zap.Error(servErr))
		}
	}()

	s.log.Info("server started", zap.String("addr", addr))

	return nil
}

func (s *Service) Get(ctx context.Context, req *proto2.CacheGetRequest) (*proto2.CacheGetResponse, error) {
	url, err := s.cacheService.Get(ctx, req.ShortCode)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "key not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get cache: %v", err)
	}
	return &proto2.CacheGetResponse{Url: url}, err
}

func (s *Service) Set(ctx context.Context, req *proto2.CacheSetRequest) (*emptypb.Empty, error) {
	if req.ShortCode == "" || req.Url == "" {
		return nil, status.Errorf(codes.InvalidArgument, "short code and url must not be empty")
	}
	err := s.cacheService.Set(ctx, req.ShortCode, req.Url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set cache: %v", err)
	}
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, req *proto2.CacheDeleteRequest) (*emptypb.Empty, error) {
	err := s.cacheService.Del(ctx, req.ShortCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete cache: %v", err)
	}
	return nil, nil

}

func (s *Service) Shutdown() {
	s.grpcServer.Stop()
	s.log.Debug("grpc shut downed")
}
