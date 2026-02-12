package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/model"
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/util"
	"github.com/gin-gonic/gin"
)

type ClickEvent struct {
	ShortCode string
	ClickedAt time.Time
	IP        string
	UserAgent string
	Referer   string
	Country   string
}

func (e *ClickEvent) WriteClickEvent(c *gin.Context) {
	e.ShortCode = c.Param("code")
	e.ClickedAt = time.Now()
	//e.IP = c.ClientIP()
	e.UserAgent = c.GetHeader("User-Agent")
	e.Referer = c.GetHeader("Referer")
	ip, err := util.CountryByIP(c.ClientIP())
	if err != nil {
		e.Country = "e"
	}
	e.IP = ip

}

func (e *ClickEvent) BindToRequest() *model.ClickRequest {
	return &model.ClickRequest{
		ShortCode: e.ShortCode,
		ClickedAt: e.ClickedAt,
		IP:        e.IP,
		UserAgent: e.UserAgent,
		Referer:   e.Referer,
		Country:   e.Country,
	}
}
