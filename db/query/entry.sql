-- name: CreateEntry :one
INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE account_id = $1;
