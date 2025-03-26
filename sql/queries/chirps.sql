-- name: CreateChirp :one
INSERT INTO
  chirps (id, created_at, updated_at, body, user_id)
VALUES
  (gen_random_uuid (), NOW (), NOW (), $1, $2) RETURNING *;

-- name: GetAllChiprs :many
SELECT
  *
FROM
  chirps;

-- name: GetChirpsByAuthorId :many
SELECT
  *
FROM
  chirps
WHERE user_id = @author_id::uuid;


-- name: DeleteAllChrips :exec
DELETE FROM
  chirps;


-- name: DeleteChripById :exec
DELETE FROM
  chirps
WHERE
  id = $1;


-- name: GetChirpById :one
SELECT
  *
FROM
  chirps
WHERE
  id = $1;
