-- name: ListPersons :many
SELECT id, name, birth_date FROM persons;

-- name: CreatePerson :one
INSERT INTO persons (id, name, birth_date)
VALUES (SUBSTRING(encode(sha256(($1 || $2 || NOW()::text)::bytea), 'hex'), 1, 16), $1, $2)
RETURNING *;

-- name: DeletePerson :exec
DELETE FROM persons WHERE id = $1;