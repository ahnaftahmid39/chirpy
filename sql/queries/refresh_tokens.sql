-- name: CreateRefreshToken :one
INSERT INTO
  refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
  )
VALUES
  ($1, NOW (), NOW (), $2, $3, NULL) RETURNING *;


-- name: DeleteAllRefreshTokens :exec
DELETE FROM
  refresh_tokens;


-- name: GetRefreshTokenByUserId :one
SELECT
  *
FROM
  refresh_tokens
WHERE
  user_id = $1
LIMIT
  1;


-- name: GetRefreshTokenByToken :one
SELECT
  *
FROM
  refresh_tokens
WHERE
  token = $1
LIMIT
  1;


-- name: RevokeRefreshTokenByToken :one
UPDATE
  refresh_tokens
SET
  revoked_at = $2,
  updated_at = $3
WHERE
  token = $1
RETURNING *;
