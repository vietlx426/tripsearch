-- name: InsertUser :one
INSERT INTO users (email, password, full_name, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: FindUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;
