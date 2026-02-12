package entity

import "time"

type Click struct {
	Id            int       `db:"id"`
	ShortCode     string    `db:"short_code"`
	LastClickedAt time.Time `db:"clicked_at"`
	IP            string    `db:"ip_address"`
	UserAgent     string    `db:"user_agent"`
	Referer       string    `db:"referer"`
	Country       string    `db:"country"`
}
