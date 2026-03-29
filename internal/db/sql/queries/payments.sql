-- name: GetPaymentByID :one
SELECT id, user_id, amount, currency
FROM payments
WHERE id = $1;

