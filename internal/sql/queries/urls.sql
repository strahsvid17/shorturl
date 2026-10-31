-- name: CreateURLAlias :one
INSERT INTO
  urls (url, alias)
VALUES
  ($1, $2)
RETURNING
  url_id;
