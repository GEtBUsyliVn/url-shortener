-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS clicks
(
    id         BIGSERIAL PRIMARY KEY,

    short_code VARCHAR(10) NOT NULL,

    clicked_at TIMESTAMP   NOT NULL DEFAULT NOW(),

    ip_address INET        NULL,

    user_agent TEXT        NULL,

    referer    TEXT        NULL,

    country    VARCHAR(2)  NULL
);

CREATE TABLE IF NOT EXISTS url_stats
(
    short_code      VARCHAR(10) PRIMARY KEY,

    total_clicks    BIGINT    NOT NULL DEFAULT 0,

    unique_visitors BIGINT    NOT NULL DEFAULT 0,

    last_clicked_at TIMESTAMP NULL,

    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для clicks
CREATE INDEX IF NOT EXISTS idx_clicks_short_code
    ON clicks(short_code);

CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at
    ON clicks(clicked_at);

CREATE INDEX IF NOT EXISTS idx_clicks_ip_address
    ON clicks(ip_address);

-- Индекс для url_stats
CREATE INDEX IF NOT EXISTS idx_url_stats_updated_at
    ON url_stats(updated_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_clicks_short_code;
DROP INDEX IF EXISTS idx_clicks_clicked_at;
DROP INDEX IF EXISTS idx_clicks_ip_address;
DROP INDEX IF EXISTS idx_url_stats_updated_at;

DROP TABLE IF EXISTS clicks;
DROP TABLE IF EXISTS url_stats;

-- +goose StatementEnd