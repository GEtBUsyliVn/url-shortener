package postgres

import (
	"context"

	"github.com/GEtBUsyliVn/url-shortener/url/repository/entity"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type BasicRepository struct {
	db *sqlx.DB
	//log *.Logger
	log *zap.Logger
}

func NewRepository(log *zap.Logger, db *sqlx.DB) *BasicRepository {
	return &BasicRepository{
		db:  db,
		log: log.Named("storage"),
	}
}

func (r *BasicRepository) Create(ctx context.Context, url *entity.Url) error {
	q := `INSERT INTO url (short_code, original_url, created_at,expires_at,user_id,is_active)
		  VALUES (:short_code, :original_url, :created_at,:expires_at,:user_id,:is_active)`
	_, err := r.db.NamedExecContext(ctx, q, url)
	return err
}

func (r *BasicRepository) Get(ctx context.Context, shortCode string) (*entity.Url, error) {
	q := `SELECT id, short_code, original_url, created_at,expires_at,user_id,is_active
		  FROM url
		  WHERE short_code = $1`
	var url entity.Url
	err := r.db.GetContext(ctx, &url, q, shortCode)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *BasicRepository) Delete(ctx context.Context, shortCode string) error {
	q := `DELETE FROM url WHERE short_code = $1`
	_, err := r.db.ExecContext(ctx, q, shortCode)
	return err
}

func (r *BasicRepository) Exists(ctx context.Context, code string) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM url WHERE short_code = $1)`,
		code,
	).Scan(&exists)

	return exists, err
}
