-- +goose Up
CREATE TABLE IF NOT EXISTS filters (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    query       JSONB       NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS filters;