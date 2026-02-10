package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/analytics/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func NewDB(config *config.Config, log *zap.Logger) (*sqlx.DB, error) {
	c := config.Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DataBase)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`SET TIME ZONE 'UTC'`)
	if err != nil {
		return nil, err
	}

	// проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	var exists bool
	_ = db.QueryRowContext(ctx, `
    select exists(
        select 1 from information_schema.tables
        where table_schema='public' and table_name='url'
    )
`).Scan(&exists)
	//if err != nil {  }
	if !exists {
		log.Fatal("url table not found in connected database")
	}
	return db, nil
}
