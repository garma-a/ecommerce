-- name: GetProduct :one
-- Returns a single struct. Errors if ID is not found.
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
-- Returns a slice of structs. Returns empty slice if no products exist.
SELECT * FROM products ORDER BY id;

-- name: CreateProduct :one
INSERT INTO products (name, price)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
