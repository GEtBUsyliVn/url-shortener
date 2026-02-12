package model

import "time"

type MemoryCache struct {
	Url       string
	ExpiresAt time.Time
}
