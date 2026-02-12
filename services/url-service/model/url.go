package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/url-service/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/repository/entity"
)

type Url struct {
	Id          int
	ShortCode   string
	OriginalUrl string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	UserId      string
	IsActive    bool
}

func (u *Url) Entity() *entity.Url {
	return &entity.Url{
		Id:          u.Id,
		ShortCode:   u.ShortCode,
		OriginalUrl: u.OriginalUrl,
		CreatedAt:   u.CreatedAt.UTC().Truncate(time.Second),
		ExpiresAt:   u.ExpiresAt,
		UserId:      u.UserId,
		IsActive:    u.IsActive,
	}
}

func (u *Url) BindProtoRequest(p *proto.CreateURLRequest) {
	u.OriginalUrl = p.GetOriginalUrl()
	u.UserId = p.GetUserId()
	u.ExpiresAt = p.ExpiresAt.AsTime()
}
