package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	proto2 "github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type Client struct {
	log        *zap.Logger
	address    string
	m          sync.Mutex
	client     proto2.URLServiceClient
	tlsEnabled bool
}

func NewGrpcClient(address string, tlsEnabled bool, log *zap.Logger) *Client {
	return &Client{
		log:        log.Named("grpc client"),
		address:    address,
		tlsEnabled: tlsEnabled,
	}
}

// url/pkg/api/grpc/proto/schema/url.proto
func (c *Client) getClient() (proto2.URLServiceClient, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.client != nil {
		return c.client, nil
	}

	creds := insecure.NewCredentials()
	var err error

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                defaultKeepaliveTime,
				Timeout:             defaultKeepaliveTimeout,
				PermitWithoutStream: true,
			},
		),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultClientMaxReceiveMessageSize),
			grpc.MaxCallSendMsgSize(defaultClientMaxSendMessageSize),
		),
	}

	conn, err := grpc.NewClient(c.address, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	c.client = proto2.NewURLServiceClient(conn)

	return c.client, nil
}

func (c *Client) CreateUrl(ctx context.Context, req *model.CreateUrlRequest) (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	resp, err := client.CreateShortURL(ctx, req.Proto())
	if err != nil {
		return "", fmt.Errorf("create short url: %w", err)
	}

	return resp.ShortCode, nil
}

func (c *Client) GetOriginalUrl(ctx context.Context, shortCode string) (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	url, err := client.GetOriginalURL(ctx, &proto2.GetURLRequest{ShortCode: shortCode})
	if err != nil {
		if errors.Is(err, status.Error(codes.NotFound, "url not found")) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("get original url: %w", err)
	}

	return url.OriginalUrl, nil

}
