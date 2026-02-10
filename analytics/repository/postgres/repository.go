package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GEtBUsyliVn/url-shortener/analytics/repository"
	"github.com/GEtBUsyliVn/url-shortener/analytics/repository/entity"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type BasicRepository struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewRepository(log *zap.Logger, db *sqlx.DB) *BasicRepository {
	return &BasicRepository{
		db:  db,
		log: log.Named("storage"),
	}
}

func (r *BasicRepository) CreateClick(ctx context.Context, click *entity.Click) error {
	q := `INSERT INTO clicks (short_code,clicked_at, ip_address, user_agent, referer, country )
		  VALUES (:short_code,:clicked_at, :ip_address, :user_agent, :referer, :country)`
	_, err := r.db.NamedExecContext(ctx, q, click)
	return err
}

func (r *BasicRepository) CreateStats(ctx context.Context, shortCode string) error {
	q := `INSERT INTO url_stats (short_code, total_clicks, unique_visitors, last_clicked_at, updated_at)
		  SELECT
          $1::varchar(10)                              AS short_code,
  		  COUNT(*)                                     AS total_clicks,
  		  COUNT(DISTINCT c.ip_address)                 AS unique_visitors,
  		  MAX(c.clicked_at)                            AS last_clicked_at,
  		  NOW()                                        AS updated_at
		  FROM clicks c
		  WHERE c.short_code = $1
		  ON CONFLICT (short_code) DO UPDATE
		  SET total_clicks    = EXCLUDED.total_clicks,
		      unique_visitors = EXCLUDED.unique_visitors,
		      last_clicked_at = EXCLUDED.last_clicked_at,
		      updated_at      = EXCLUDED.updated_at;`

	_, err := r.db.ExecContext(ctx, q, shortCode)
	return err
}

func (r *BasicRepository) GetStatistics(ctx context.Context, shortCode string) (*entity.Statistics, error) {
	q := `SELECT short_code, total_clicks, unique_visitors, last_clicked_at, updated_at
		  FROM url_stats
		  WHERE short_code = $1`

	var stats entity.Statistics
	err := r.db.GetContext(ctx, &stats, q, shortCode)
	if errors.Is(err, sql.ErrNoRows) {
		r.log.Error("failed to get statistics", zap.String("short_code", shortCode), zap.Error(err))
		return nil, repository.ErrNotFound
	}
	return &stats, nil
}
