package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "sub-hub/internal/domain/errors"
	"sub-hub/internal/domain/model"
	"sub-hub/internal/domain/ports"
)

type SubscriptionsRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptions(pool *pgxpool.Pool) *SubscriptionsRepo {
	return &SubscriptionsRepo{pool: pool}
}

func (r *SubscriptionsRepo) Create(ctx context.Context, sub model.Subscription) (model.Subscription, error) {
	q := `
INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
`

	row := r.pool.QueryRow(ctx, q, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	var out model.Subscription
	if err := row.Scan(&out.ID, &out.ServiceName, &out.Price, &out.UserID, &out.StartDate, &out.EndDate, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return model.Subscription{}, mapPGError(err)
	}
	return out, nil
}

func (r *SubscriptionsRepo) Get(ctx context.Context, id string) (model.Subscription, error) {
	q := `
SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
FROM subscriptions
WHERE id = $1
`
	row := r.pool.QueryRow(ctx, q, id)
	var out model.Subscription
	if err := row.Scan(&out.ID, &out.ServiceName, &out.Price, &out.UserID, &out.StartDate, &out.EndDate, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, domainerrors.ErrNotFound
		}
		return model.Subscription{}, mapPGError(err)
	}
	return out, nil
}

func (r *SubscriptionsRepo) Update(ctx context.Context, sub model.Subscription) (model.Subscription, error) {
	q := `
UPDATE subscriptions
SET service_name = $2,
    price = $3,
    user_id = $4,
    start_date = $5,
    end_date = $6,
    updated_at = now()
WHERE id = $1
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
`
	row := r.pool.QueryRow(ctx, q, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	var out model.Subscription
	if err := row.Scan(&out.ID, &out.ServiceName, &out.Price, &out.UserID, &out.StartDate, &out.EndDate, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, domainerrors.ErrNotFound
		}
		return model.Subscription{}, mapPGError(err)
	}
	return out, nil
}

func (r *SubscriptionsRepo) Delete(ctx context.Context, id string) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return mapPGError(err)
	}
	if ct.RowsAffected() == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func (r *SubscriptionsRepo) List(ctx context.Context, f ports.SubscriptionListFilter) ([]model.Subscription, error) {
	limit := f.Limit
	if limit == 0 {
		limit = 50
	}
	q := `
SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
FROM subscriptions
WHERE ($1::uuid IS NULL OR user_id = $1::uuid)
  AND ($2::text IS NULL OR service_name = $2::text)
  AND ($3::date IS NULL OR coalesce(end_date, $3::date) >= $3::date)
  AND ($4::date IS NULL OR start_date <= $4::date)
ORDER BY start_date DESC, created_at DESC
LIMIT $5 OFFSET $6
`

	var userID any = nil
	if f.UserID != nil {
		userID = *f.UserID
	}
	var service any = nil
	if f.ServiceName != nil {
		service = *f.ServiceName
	}
	var from any = nil
	if f.From != nil {
		from = *f.From
	}
	var to any = nil
	if f.To != nil {
		to = *f.To
	}

	rows, err := r.pool.Query(ctx, q, userID, service, from, to, limit, f.Offset)
	if err != nil {
		return nil, mapPGError(err)
	}
	defer rows.Close()

	var out []model.Subscription
	for rows.Next() {
		var s model.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, mapPGError(err)
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, mapPGError(err)
	}
	return out, nil
}

func (r *SubscriptionsRepo) TotalCost(ctx context.Context, f ports.SubscriptionTotalFilter) (int64, error) {
	q := `
SELECT coalesce(sum(
  CASE
    WHEN overlap_end < overlap_start THEN 0
    ELSE (
      (
        (extract(year from overlap_end)::int * 12 + extract(month from overlap_end)::int)
        - (extract(year from overlap_start)::int * 12 + extract(month from overlap_start)::int)
        + 1
      ) * price
    )
  END
), 0)::bigint AS total
FROM (
  SELECT price,
         greatest(start_date, $1::date) AS overlap_start,
         least(coalesce(end_date, $2::date), $2::date) AS overlap_end
  FROM subscriptions
  WHERE start_date <= $2::date
    AND coalesce(end_date, $2::date) >= $1::date
    AND ($3::uuid IS NULL OR user_id = $3::uuid)
    AND ($4::text IS NULL OR service_name = $4::text)
) s
`

	var userID any = nil
	if f.UserID != nil {
		userID = *f.UserID
	}
	var service any = nil
	if f.ServiceName != nil {
		service = *f.ServiceName
	}

	var total int64
	if err := r.pool.QueryRow(ctx, q, f.From, f.To, userID, service).Scan(&total); err != nil {
		return 0, mapPGError(err)
	}
	return total, nil
}

func mapPGError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return domainerrors.ErrConflict
		}
		return fmt.Errorf("db error: %s", pgErr.Code)
	}
	return err
}
