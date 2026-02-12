package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/grpc/proto"
	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository/entity"
)

type Click struct {
	Id        int
	ShortCode string
	ClickedAt time.Time
	IP        string
	UserAgent string
	Referer   string
	Country   string
}

func (c *Click) Entity() *entity.Click {
	return &entity.Click{
		Id:            c.Id,
		ShortCode:     c.ShortCode,
		LastClickedAt: c.ClickedAt,
		IP:            c.IP,
		UserAgent:     c.UserAgent,
		Referer:       c.Referer,
		Country:       c.Country,
	}
}

func (c *Click) BindProto(req *proto.ClickEvent) {
	c.ShortCode = req.ShortCode
	c.ClickedAt = req.ClickedAt.AsTime().UTC()
	c.IP = req.IpAddress
	c.UserAgent = req.UserAgent
	c.Referer = req.Referer
	c.Country = req.Country
}
