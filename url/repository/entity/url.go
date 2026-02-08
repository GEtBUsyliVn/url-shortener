package entity

import (
	"database/sql"
	"time"
)

type Url struct {
	Id          int          `db:"id"`
	ShortCode   string       `db:"short_code"`
	OriginalUrl string       `db:"original_url"`
	CreatedAt   time.Time    `db:"created_at"`
	ExpiresAt   sql.NullTime `db:"expires_at"`
	UserId      string       `db:"user_id"`
	IsActive    bool         `db:"is_active"`
}
