package grpc

import (
	"context"
	"fmt"
	"sync"

	proto2 "github.com/GEtBUsyliVn/url-shortener/services/cache-service/pkg/api/grpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	log        *zap.Logger
	address    string
	m          sync.Mutex
	client     proto2.CacheServiceClient
	tlsEnabled bool
}

func NewGrpcClient(address string, tlsEnabled bool, log *zap.Logger) *Client {
	return &Client{
		log:        log.Named("grpc client"),
		address:    address,
		tlsEnabled: tlsEnabled,
	}
}

func (c *Client) getClient() (proto2.CacheServiceClient, error) {
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

	c.client = proto2.NewCacheServiceClient(conn)

	return c.client, nil
}

func (c *Client) Get(ctx context.Context, code string) (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	resp, err := client.Get(ctx, &proto2.CacheGetRequest{ShortCode: code})
	if err != nil {
		return "", fmt.Errorf("failed to get cache: %w", err)
	}

	return resp.Url, nil
}

func (c *Client) Set(ctx context.Context, code, url string) error {
	client, err := c.getClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	_, err = client.Set(ctx, &proto2.CacheSetRequest{ShortCode: code, Url: url})
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (c *Client) Del(ctx context.Context, code string) error {
	client, err := c.getClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	_, err = client.Delete(ctx, &proto2.CacheDeleteRequest{ShortCode: code})
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}
