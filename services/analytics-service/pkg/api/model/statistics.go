package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc/proto"
)

type StatsRequest struct {
	ShortCode string
}

func (s *StatsRequest) Proto() *proto.StatsRequest {
	return &proto.StatsRequest{
		ShortCode: s.ShortCode,
	}
}

type StatsResponse struct {
	ShortCode    string
	TotalClicks  int64
	UniqVisitors int64
	ClickedAt    time.Time
}

func (s *StatsResponse) Proto(resp *proto.StatsResponse) {
	s.ShortCode = resp.ShortCode
	s.TotalClicks = resp.TotalClicks
	s.UniqVisitors = resp.UniqueVisitors
	s.ClickedAt = resp.ClickedAt.AsTime()
}
