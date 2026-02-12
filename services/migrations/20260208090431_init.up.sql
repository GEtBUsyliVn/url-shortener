
CREATE TABLE IF NOT EXISTS clicks
(
    id         BIGSERIAL PRIMARY KEY,

    short_code VARCHAR(10) NOT NULL,

    clicked_at TIMESTAMP   NOT NULL DEFAULT NOW(),

    ip_address TEXT        NULL,

    user_agent TEXT        NULL,

    referer    TEXT        NULL,

    country    VARCHAR(10)  NULL
);

CREATE TABLE IF NOT EXISTS url_stats
(
    short_code      VARCHAR(10) PRIMARY KEY,

    total_clicks    BIGINT    NOT NULL DEFAULT 0,

    unique_visitors BIGINT    NOT NULL DEFAULT 0,

    last_clicked_at TIMESTAMP NULL,

    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_clicks_short_code
    ON clicks(short_code);

CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at
    ON clicks(clicked_at);

CREATE INDEX IF NOT EXISTS idx_clicks_ip_address
    ON clicks(ip_address);


CREATE INDEX IF NOT EXISTS idx_url_stats_updated_at
    ON url_stats(updated_at);

CREATE TABLE IF NOT EXISTS url
(
    id           BIGSERIAL PRIMARY KEY,

    short_code   VARCHAR(10) UNIQUE NOT NULL,

    original_url TEXT               NOT NULL,

    created_at   TIMESTAMP DEFAULT NOW(),

    expires_at   TIMESTAMP          NULL,

    user_id      VARCHAR(50)        NULL,

    is_active    BOOLEAN   DEFAULT true
);