package model

import (
	"database/sql"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/url/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/url/repository/entity"
)

type Url struct {
	Id          int
	ShortCode   string
	OriginalUrl string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	UserId      string
	IsActive    bool
}

func (u *Url) Entity() *entity.Url {
	var exp sql.NullTime
	if u.ExpiresAt != nil {
		exp = sql.NullTime{Time: u.ExpiresAt.UTC().Truncate(time.Second), Valid: true}
	}

	return &entity.Url{
		Id:          u.Id,
		ShortCode:   u.ShortCode,
		OriginalUrl: u.OriginalUrl,
		CreatedAt:   u.CreatedAt.UTC().Truncate(time.Second),
		ExpiresAt:   exp,
		UserId:      u.UserId,
		IsActive:    u.IsActive,
	}
}

func (u *Url) BindProtoRequest(p *proto.CreateURLRequest) error {
	u.OriginalUrl = p.GetOriginalUrl()
	u.UserId = p.GetUserId()
	u.ShortCode = p.GetShortCode()

	ts := p.GetExpiresAt() //

	if ts == nil {
		u.ExpiresAt = nil
		return nil
	}

	if err := ts.CheckValid(); err != nil {
		return err
	}

	t := ts.AsTime().UTC().Truncate(time.Second)
	u.ExpiresAt = &t

	return nil
}
