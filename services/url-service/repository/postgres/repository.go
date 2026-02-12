package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GEtBUsyliVn/url-shortener/services/url-service/repository"
	"github.com/GEtBUsyliVn/url-shortener/services/url-service/repository/entity"
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
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &url, err
}

//func (r *BasicRepository) Update(ctx context.Context, url *entity.Url) error {
//	q := `UPDATE "url"
//		  SET original_url = :original_url,
//		      created_at   = :created_at,
//		      expires_at   = :expires_at,
//		      user_id      = :user_id,
//		      is_active    = :is_active
//		  WHERE short_code  = :short_code`
//
//	res, err := r.db.NamedExecContext(ctx, q, url)
//	if err != nil {
//		return err
//	}
//	n, err := res.RowsAffected()
//	if err != nil {
//		return err
//	}
//	if n == 0 {
//		return fmt.Errorf("update: no rows affected (short_code=%s)", url.ShortCode)
//	}
//	return nil
//}

func (r *BasicRepository) UpdateExpired(ctx context.Context) (int, error) {
	q := `UPDATE url SET is_active = false WHERE expires_at < NOW() AND is_active = true`
	res, err := r.db.ExecContext(ctx, q)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func (r *BasicRepository) List(ctx context.Context) ([]*entity.Url, error) {
	q := `SELECT id, short_code, original_url, created_at,expires_at,user_id,is_active FROM url`
	var urls []*entity.Url
	err := r.db.SelectContext(ctx, &urls, q)
	return urls, err
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
