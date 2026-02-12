package model

import "time"

type StatsRest struct {
	ShortCode    string    `json:"shortCode"`
	TotalClicks  int64     `json:"totalClicks"`
	UniqVisitors int64     `json:"uniqVisitors"`
	ClickedAt    time.Time `json:"clickedAt"`
}
