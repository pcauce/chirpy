-- name: StoreRefresh :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
        $1,
        now(),
        now(),
        $2,
        $3
       );

-- name: RevokeRefresh :exec
UPDATE refresh_tokens
SET updated_at = $2, revoked_at = $2
WHERE token = $1;

-- name: GetUserFromRefresh :one
-- name: GetUserFromRefresh :one
SELECT user_id
FROM refresh_tokens
WHERE token = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();