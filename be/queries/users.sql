-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = ?;

-- name: CreateUser :exec
INSERT INTO
    users (name, email)
VALUES
    (?, ?);

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = ?;