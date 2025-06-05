-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
           gen_random_uuid(),
           now(),
           now(),
           $1,
           $2
       )
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeChirpyRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;