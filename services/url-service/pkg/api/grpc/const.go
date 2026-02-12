package grpc

import (
	"math"
	"time"
)

const (
	defaultKeepaliveTime               = 60 * time.Second
	defaultKeepaliveTimeout            = 10 * time.Second
	defaultClientMaxReceiveMessageSize = math.MaxInt64
	defaultClientMaxSendMessageSize    = 10 * 1024 * 1024
)
