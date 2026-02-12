package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ClickRequest struct {
	ShortCode string
	ClickedAt time.Time
	IP        string
	UserAgent string
	Referer   string
	Country   string
}

func (r *ClickRequest) Proto() *proto.ClickEvent {
	return &proto.ClickEvent{
		ShortCode: r.ShortCode,
		ClickedAt: timestamppb.New(r.ClickedAt),
		IpAddress: r.IP,
		UserAgent: r.UserAgent,
		Referer:   r.Referer,
		Country:   r.Country,
	}
}
