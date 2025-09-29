-- name: CreateTicket :execresult
INSERT INTO tickets (title, description, created_by, priority,status)
VALUES (?, ?, ?, ?,?);


-- name: ListTickets :many
SELECT
    id,
    title,
    description,
    status,
    priority,
    created_by,
    assigned_to,
    created_at,
    updated_at
FROM tickets
ORDER BY created_at DESC
LIMIT ? OFFSET ?;


-- name: GetTicket :one
SELECT * FROM tickets
WHERE id = ? LIMIT 1;

-- name: UpdateTicketStatus :exec
UPDATE tickets
SET status = ?, updated_at = NOW()
WHERE id = ?;

-- name: AssignTicket :exec
UPDATE tickets
SET assigned_to = ?, updated_at = NOW()
WHERE id = ?;
-- name: GetTicketByTitleAndUser :one
SELECT * FROM tickets
WHERE title = ? AND created_by = ?
LIMIT 1;


-- name: CreateUser :execresult
INSERT INTO users (full_name, email,  role)
VALUES (?, ?, ?);
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT
    id,
    full_name,
    email,
    role,
    created_at
FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateUser :exec
UPDATE users
SET
  full_name = sqlc.arg(full_name),
  email     = sqlc.arg(email),
  role      = sqlc.arg(role)
WHERE id = sqlc.arg(id);

-- name: GetUserByEmailExcludingID :one
SELECT * FROM users
WHERE email = ? AND id != ?
LIMIT 1;

-- name: CreateTransaction :execresult
INSERT INTO transactions (transaction_id, user_id,amount,currency,status,payment_method)
VALUES(?,?,?,?,?,?);

-- name: ListTransactions :many
SELECT
id,
transaction_id,
amount,
currency,
payment_method
FROM transactions
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetTransanctionByID :one
SELECT *FROM transactions
WHERE id = ?
LIMIT 1;

-- name: CreateCustomer :execresult
INSERT INTO customers (full_name, email, phone_number)
VALUES (?, ?, ?);

-- name: GetCustomerByEmail :one
SELECT * FROM customers
WHERE email = ? LIMIT 1;
