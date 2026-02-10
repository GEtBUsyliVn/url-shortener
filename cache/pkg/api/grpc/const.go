package grpc

import (
	"math"
	"time"
)

const (
	defaultKeepaliveTime               = 10 * time.Second
	defaultKeepaliveTimeout            = time.Second
	defaultClientMaxReceiveMessageSize = math.MaxInt64
	defaultClientMaxSendMessageSize    = 10 * 1024 * 1024
	defaultContextTimeout              = 30 * time.Second
)
