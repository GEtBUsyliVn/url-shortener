package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/GEtBUsyliVn/url-shortener/services/url-service/config"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/model"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/grpc/proto"
	service2 "github.com/GEtBUsyliVn/url-shortener/services/url-service/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	urlService *service2.Service
	proto.UnimplementedURLServiceServer
	listener net.Listener
}

func NewGrpcService(log *zap.Logger, urlService *service2.Service) *Service {
	return &Service{
		log:        log.Named("grpc url"),
		urlService: urlService,
		grpcServer: grpc.NewServer(),
	}
}

func (s *Service) Init(cfg config.GRPC) error {
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("unable to listen port %w", err)
	}
	s.listener = listener

	proto.RegisterURLServiceServer(s.grpcServer, s)

	go func() {
		s.log.Debug("grpc serve", zap.String("address", cfg.Addr))

		if servErr := s.grpcServer.Serve(s.listener); servErr != nil {
			s.log.Panic("unable to serve", zap.Error(servErr))
		}
	}()

	s.log.Info("server started", zap.String("addr", cfg.Addr))

	return nil
}

func (s *Service) CreateShortURL(ctx context.Context, req *proto.CreateURLRequest) (*proto.CreateURLResponse, error) {
	if req.OriginalUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}
	url := &model.Url{}
	url.BindProtoRequest(req)
	res, err := s.urlService.CreateShortURL(ctx, url)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "create short url: %v", err)
	}

	return &proto.CreateURLResponse{ShortCode: res}, nil
}

func (s *Service) GetOriginalURL(ctx context.Context, req *proto.GetURLRequest) (*proto.GetURLResponse, error) {
	if req.ShortCode == "" {
		return nil, status.Error(codes.InvalidArgument, "short code is required")
	}

	res, err := s.urlService.GetShortUrl(ctx, req.ShortCode)
	if err != nil {
		if errors.Is(err, service2.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "url not found")
		}
		return nil, status.Errorf(codes.Internal, "get original url: %v", err)
	}

	return &proto.GetURLResponse{OriginalUrl: res}, nil
}

func (s *Service) Shutdown() {
	s.grpcServer.Stop()
	s.log.Debug("grpc shut downed")
}
