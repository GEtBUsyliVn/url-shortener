-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd
