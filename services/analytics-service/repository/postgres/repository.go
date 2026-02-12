package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository"
	entity2 "github.com/GEtBUsyliVn/url-shortener/services/analytics-service/repository/entity"
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

func (r *BasicRepository) CreateClick(ctx context.Context, click *entity2.Click) error {
	q := `INSERT INTO clicks (short_code,clicked_at, ip_address, user_agent, referer, country )
		  VALUES (:short_code,:clicked_at, :ip_address, :user_agent, :referer, :country)`
	_, err := r.db.NamedExecContext(ctx, q, click)
	r.log.Info("storage error", zap.Error(err))
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

func (r *BasicRepository) GetStatistics(ctx context.Context, shortCode string) (*entity2.Statistics, error) {
	q := `SELECT short_code, total_clicks, unique_visitors, last_clicked_at, updated_at
		  FROM url_stats
		  WHERE short_code = $1`

	var stats entity2.Statistics
	err := r.db.GetContext(ctx, &stats, q, shortCode)
	if errors.Is(err, sql.ErrNoRows) {
		r.log.Error("failed to get statistics", zap.String("short_code", shortCode), zap.Error(err))
		return nil, repository.ErrNotFound
	}
	return &stats, nil
}

func (r *BasicRepository) GetUniqClicks(ctx context.Context) ([]*entity2.Click, error) {
	q := `SELECT DISTINCT ON(short_code) id,short_code, clicked_at, ip_address, user_agent, referer, country
		FROM clicks`
	var click []*entity2.Click
	err := r.db.SelectContext(ctx, &click, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Info("no clicks found for aggregation")
			return nil, nil
		}
		r.log.Error("failed to get unique clicks", zap.Error(err))
		return nil, err
	}
	return click, nil
}

func (r *BasicRepository) CreateClicksBatch(ctx context.Context, clicks []*entity2.Click) error {
	if len(clicks) == 0 {
		return nil
	}

	var (
		sb   strings.Builder
		args = make([]any, 0, len(clicks)*4)
	)

	sb.WriteString(
		`
		INSERT INTO clicks (short_code, clicked_at, ip_address, user_agent,referer,country)
		VALUES`)

	argPos := 1

	for i, c := range clicks {
		if i > 0 {
			sb.WriteString(",")
		}

		sb.WriteString(fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d)",
			argPos,
			argPos+1,
			argPos+2,
			argPos+3,
			argPos+4,
			argPos+5,
		))

		args = append(
			args,
			c.ShortCode,
			c.LastClickedAt,
			c.IP,
			c.UserAgent,
			c.Referer,
			c.Country,
		)

		argPos += 4
	}

	query := sb.String()

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
