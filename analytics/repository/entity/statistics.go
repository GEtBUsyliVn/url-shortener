package entity

import "time"

type Statistics struct {
	ShortCode     string    `db:"short_code"`
	TotalClicks   int       `db:"total_clicks"`
	UniqVisitors  int       `db:"unique_visitors"`
	LastClickedAt time.Time `db:"last_clicked_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
