-- name: GetUserIDFromIndex :one
SELECT id
FROM users
ORDER BY created_at
OFFSET $1
LIMIT 1;

-- name: GetChirpIDFromIndex :one
SELECT id
FROM chirps
ORDER BY created_at
OFFSET $1
LIMIT 1;
