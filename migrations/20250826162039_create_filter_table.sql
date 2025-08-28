-- +goose Up
CREATE TABLE IF NOT EXISTS filters (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    query       JSONB       NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_filters_updated_at ON filters;

CREATE TRIGGER trg_filters_updated_at
BEFORE UPDATE ON filters
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_filters_updated_at ON filters;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS filters;