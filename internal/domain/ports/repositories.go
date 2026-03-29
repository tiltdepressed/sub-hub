package ports

import (
	"context"
	"time"

	"sub-hub/internal/domain/model"
)

type SubscriptionsRepository interface {
	Create(ctx context.Context, sub model.Subscription) (model.Subscription, error)
	Get(ctx context.Context, id string) (model.Subscription, error)
	Update(ctx context.Context, sub model.Subscription) (model.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, f SubscriptionListFilter) ([]model.Subscription, error)
	TotalCost(ctx context.Context, f SubscriptionTotalFilter) (int64, error)
}

type SubscriptionListFilter struct {
	UserID      *string
	ServiceName *string
	From        *time.Time
	To          *time.Time
	Limit       int
	Offset      int
}

type SubscriptionTotalFilter struct {
	UserID      *string
	ServiceName *string
	From        time.Time
	To          time.Time
}
