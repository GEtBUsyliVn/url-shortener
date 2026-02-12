package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/pkg/api/model"
	modelR "github.com/GEtBUsyliVn/url-shortener/services/api-gateway/rest/model"
)

type Stats struct {
	ShortCode    string
	TotalClicks  int64
	UniqVisitors int64
	ClickedAt    time.Time
}

func (s *Stats) BindRest() *modelR.StatsRest {
	return &modelR.StatsRest{
		ShortCode:    s.ShortCode,
		TotalClicks:  s.TotalClicks,
		UniqVisitors: s.UniqVisitors,
		ClickedAt:    s.ClickedAt,
	}
}

func (s *Stats) Bind(r *model.StatsResponse) {
	s.ShortCode = r.ShortCode
	s.TotalClicks = r.TotalClicks
	s.UniqVisitors = r.UniqVisitors
	s.ClickedAt = r.ClickedAt
}
