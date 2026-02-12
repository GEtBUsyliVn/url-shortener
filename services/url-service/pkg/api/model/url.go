package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateUrlRequest struct {
	Url      string    `json:"url" binding:"required"`
	ExpireAt time.Time `json:"expireAt"`
	UserId   string    `json:"userId"`
}

func (u *CreateUrlRequest) Proto() *proto.CreateURLRequest {
	return &proto.CreateURLRequest{
		OriginalUrl: u.Url,
		ExpiresAt:   timestamppb.New(u.ExpireAt),
		UserId:      &u.UserId,
	}
}
