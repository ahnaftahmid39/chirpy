-- name: CreateUser :one
INSERT INTO
  users (id, created_at, updated_at, email, hashed_password)
VALUES
  (
    gen_random_uuid (),
    NOW (),
    NOW (),
    $1,
    $2
  ) RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUserById :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = $4
WHERE id = $1
RETURNING *;

-- name: MakeUserChirpyRedById :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;
