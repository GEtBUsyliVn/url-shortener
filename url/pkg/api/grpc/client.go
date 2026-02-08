package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/url/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/url/pkg/api/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	log        *zap.Logger
	address    string
	m          sync.Mutex
	client     proto.URLServiceClient
	tlsEnabled bool
}

func NewGrpcClient(address string, tlsEnabled bool, log *zap.Logger) *Client {
	return &Client{
		log:        log.Named("grpc client"),
		address:    address,
		tlsEnabled: tlsEnabled,
	}
}

func (c *Client) getClient() (proto.URLServiceClient, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.client != nil {
		return c.client, nil
	}

	creds := insecure.NewCredentials()
	var err error

	//if c.tlsEnabled {
	//	var certPool *x509.CertPool
	//	if len(c.serverCA) > 0 {
	//		if certPool = x509.NewCertPool(); !certPool.AppendCertsFromPEM(c.serverCA) {
	//			return nil, fmt.Errorf("failed to add server CA's certificate")
	//		}
	//	} else if certPool, err = x509.SystemCertPool(); err != nil {
	//		return nil, fmt.Errorf("failed to get system cert pool: %w", err)
	//	}
	//
	//	creds = credentials.NewClientTLSFromCert(certPool, "")
	//}

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                defaultKeepaliveTime,
				Timeout:             defaultKeepaliveTimeout,
				PermitWithoutStream: true,
			},
		),
		//grpc.WithChainUnaryInterceptor(
		//	retry.UnaryClientInterceptor(
		//		retry.WithMax(defaultMaxRetryCount),
		//		retry.WithBackoff(retry.BackoffLinear(time.Second)),
		//		retry.WithCodes(codes.Aborted, codes.Unavailable),
		//	),
		//),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultClientMaxReceiveMessageSize),
			grpc.MaxCallSendMsgSize(defaultClientMaxSendMessageSize),
		),
	}

	conn, err := grpc.NewClient(c.address, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	c.client = proto.NewURLServiceClient(conn)

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
