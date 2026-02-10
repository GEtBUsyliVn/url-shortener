package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	proto "github.com/GEtBUsyliVn/url-shortener/cache/pkg/api/grpc/proto"
	service2 "github.com/GEtBUsyliVn/url-shortener/cache/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	log          *zap.Logger
	grpcServer   *grpc.Server
	cacheService *service2.CacheService
	proto.UnimplementedCacheServiceServer
	listener net.Listener
}

func NewGrpcService(log *zap.Logger, cacheService *service2.CacheService) *Service {
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

	proto.RegisterCacheServiceServer(s.grpcServer, s)

	go func() {
		s.log.Debug("grpc serve", zap.String("address", addr))

		if servErr := s.grpcServer.Serve(s.listener); servErr != nil {
			s.log.Panic("unable to serve", zap.Error(servErr))
		}
	}()

	s.log.Info("server started", zap.String("addr", addr))

	return nil
}

func (s *Service) Get(ctx context.Context, req *proto.CacheGetRequest) (*proto.CacheGetResponse, error) {
	url, err := s.cacheService.Get(ctx, req.ShortCode)
	if err != nil {
		if errors.Is(err, service2.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "key not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get cache: %v", err)
	}
	return &proto.CacheGetResponse{Url: url}, err
}

func (s *Service) Set(ctx context.Context, req *proto.CacheSetRequest) (*emptypb.Empty, error) {
	if req.ShortCode == "" || req.Url == "" {
		return nil, status.Errorf(codes.InvalidArgument, "short code and url must not be empty")
	}
	err := s.cacheService.Set(ctx, req.ShortCode, req.Url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set cache: %v", err)
	}
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, req *proto.CacheDeleteRequest) (*emptypb.Empty, error) {
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
