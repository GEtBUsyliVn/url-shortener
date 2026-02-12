package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Statistics struct {
	ShortCode     string
	TotalClicks   int
	UniqVisitors  int
	LastClickedAt time.Time
	UpdatedAt     time.Time
}

func (s *Statistics) Bind(entity *entity.Statistics) {
	s.ShortCode = entity.ShortCode
	s.TotalClicks = entity.TotalClicks
	s.UniqVisitors = entity.UniqVisitors
	s.LastClickedAt = entity.LastClickedAt
	s.UpdatedAt = entity.UpdatedAt
}

func (s *Statistics) BindProtoResponse() *proto.StatsResponse {
	return &proto.StatsResponse{
		ShortCode:      s.ShortCode,
		TotalClicks:    int64(s.TotalClicks),
		UniqueVisitors: int64(s.UniqVisitors),
		ClickedAt:      timestamppb.New(s.LastClickedAt),
	}
}
