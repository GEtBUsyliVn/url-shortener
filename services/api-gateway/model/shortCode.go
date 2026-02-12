package model

import (
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/rest/model"
)

type ShortCode struct {
	URL         string
	UserId      string
	ExpiredDays int
}

func (s *ShortCode) BindRestCreate(req *model.CreateShortCodeRequest) {
	s.URL = req.URL
	s.UserId = req.UserId
	s.ExpiredDays = req.ExpiredDays
}
