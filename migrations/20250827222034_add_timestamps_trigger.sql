-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_timestamps()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.created_at := now();
        NEW.updated_at := now();
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        NEW.updated_at := now();
        RETURN NEW;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_filters_updated_at ON filters;

CREATE TRIGGER trg_filters_insert_ts
BEFORE INSERT ON filters
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER trg_filters_update_ts
BEFORE UPDATE ON filters
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

DROP FUNCTION IF EXISTS set_updated_at();


-- +goose Down
DROP TRIGGER IF EXISTS trg_filters_insert_ts ON filters;
DROP TRIGGER IF EXISTS trg_filters_update_ts ON filters;

DROP FUNCTION IF EXISTS set_timestamps();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_filters_updated_at
BEFORE UPDATE ON filters
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();