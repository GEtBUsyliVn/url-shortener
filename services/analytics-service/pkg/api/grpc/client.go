package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc/proto"
	model2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/model"
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
	client     proto.AnalyticsServiceClient
	tlsEnabled bool
}

func NewGrpcClient(address string, tlsEnabled bool, log *zap.Logger) *Client {
	return &Client{
		log:        log.Named("grpc client"),
		address:    address,
		tlsEnabled: tlsEnabled,
	}
}

func (c *Client) getClient() (proto.AnalyticsServiceClient, error) {
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

	c.client = proto.NewAnalyticsServiceClient(conn)

	return c.client, nil
}

func (c *Client) ClickEvent(ctx context.Context, req *model2.ClickRequest) error {
	client, err := c.getClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	_, err = client.RecordClick(ctx, req.Proto())
	if err != nil {
		return fmt.Errorf("failed to record click: %w", err)
	}

	return nil
}

func (c *Client) GetStatistics(ctx context.Context, req *model2.StatsRequest) (*model2.StatsResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	stats, err := client.GetStatistics(ctx, req.Proto())
	if err != nil {
		if errors.Is(err, status.Errorf(codes.Internal, "failed convert to proto news: %v", err)) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	resp := &model2.StatsResponse{}
	resp.Proto(stats)
	return resp, nil
}
