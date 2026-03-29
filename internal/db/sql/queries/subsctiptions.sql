-- name: GetSubscriptionByID :one
SELECT id, user_id, plan, active
FROM subscriptions
WHERE id = $1;

