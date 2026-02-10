package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/GEtBUsyliVn/url-shortener/url/config"
	"github.com/GEtBUsyliVn/url-shortener/url/model"
	protoS "github.com/GEtBUsyliVn/url-shortener/url/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/url/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	urlService *service.Service
	protoS.UnimplementedURLServiceServer
	listener net.Listener
}

func NewGrpcService(log *zap.Logger, urlService *service.Service) *Service {
	return &Service{
		log:        log.Named("grpc url"),
		urlService: urlService,
		grpcServer: grpc.NewServer(),
	}
}

func (s *Service) Init(cfg config.GrpcConfig) error {
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("unable to listen port %w", err)
	}
	s.listener = listener

	protoS.RegisterURLServiceServer(s.grpcServer, s)

	go func() {
		s.log.Debug("grpc serve", zap.String("address", cfg.Addr))

		if servErr := s.grpcServer.Serve(s.listener); servErr != nil {
			s.log.Panic("unable to serve", zap.Error(servErr))
		}
	}()

	s.log.Info("server started", zap.String("addr", cfg.Addr))

	return nil
}

func (s *Service) CreateShortURL(ctx context.Context, req *protoS.CreateURLRequest) (*protoS.CreateURLResponse, error) {
	url := &model.Url{}
	err := url.BindProtoRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := s.urlService.CreateShortURL(ctx, url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create short url: %v", err)
	}

	return &protoS.CreateURLResponse{ShortCode: res}, nil
}

func (s *Service) Shutdown() {
	s.grpcServer.Stop()
	s.log.Debug("grpc shut downed")
}
