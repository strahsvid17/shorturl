-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
  url_id INTEGER PRIMARY KEY,
  alias TEXT NOT NULL UNIQUE,
  url TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_alias ON url_id (alias);

-- +goose Down
DROP TABLE urls CASCADE;
