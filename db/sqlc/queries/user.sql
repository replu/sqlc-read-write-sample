-- name: UserCreate :execresult
INSERT INTO users(
    name
) VALUE (
  ?
)
;

-- name: UserGet :one
SELECT *
FROM users
WHERE
  id = sqlc.arg(id)
;
