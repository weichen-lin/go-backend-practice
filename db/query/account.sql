-- name: CreateAccount :one
INSERT INTO account (owner, balance, currency) VALUES ($1, $2, $3) RETURNING id;

-- name: GetAccount :one
SELECT * FROM account WHERE id = $1;

-- name: ListAccounts :many
SELECT * FROM account 
ORDER BY last_modified_at
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account 
SET balance = $1 
WHERE id = $2
RETURNING id, balance;

-- name: DeleteAccount :exec
DELETE FROM account WHERE id = $1;