-- name: CreateTicket :execresult
INSERT INTO tickets (title, description, created_by, priority, status)
VALUES (?, ?, ?, ?, ?);

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

-- name: GetCustomers :many
SELECT * FROM customers LIMIT ? OFFSET ?;



-- db/queries.sql

-- name: GetProfileByPhone :one
SELECT id, phone, password_hash, full_name, created_at, updated_at
FROM profiles
WHERE phone = ?;

-- name: CreateOTP :execresult
INSERT INTO otp_codes (profile_id, otp_code, expires_at)
VALUES (?, ?, ?);
SELECT LAST_INSERT_ID() as id;

-- name: GetLatestOTPByProfileID :one
SELECT id, profile_id, otp_code, expires_at, verified, attempts, created_at
FROM otp_codes
WHERE profile_id = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkOTPVerified :exec
UPDATE otp_codes
SET verified = TRUE
WHERE id = ?;

-- name: IncrementOTPAttempts :exec
UPDATE otp_codes
SET attempts = attempts + 1
WHERE id = ?;

-- name: DeleteExpiredOTPs :exec
DELETE FROM otp_codes
WHERE expires_at < NOW();

-- name: CreateProfile :execresult
INSERT INTO profiles (full_name, phone, password_hash)
VALUES (?, ?, ?);

