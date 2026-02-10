package model

import (
	"time"

	"github.com/GEtBUsyliVn/url-shortener/analytics/repository/entity"
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
