-- name: CreateURLAlias :one
INSERT INTO
  urls (url, alias)
VALUES
  ($1, $2)
RETURNING
  url_id;

-- name: FindURLByAlias :one
SELECT
  url
FROM
  urls
WHERE
  alias = $1
LIMIT
  1;
