package subscriptions

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"sub-hub/internal/domain/errors"
	"sub-hub/internal/domain/model"
	"sub-hub/internal/domain/ports"
)

type Service struct {
	repo ports.SubscriptionsRepository
}

func New(repo ports.SubscriptionsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (model.Subscription, error) {
	if in.ServiceName == "" {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	if in.Price < 0 {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	if _, err := uuid.Parse(in.UserID); err != nil {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	start, err := parseMonthYear(in.StartDate)
	if err != nil {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	var end *time.Time
	if in.EndDate != nil {
		v, err := parseMonthYear(*in.EndDate)
		if err != nil {
			return model.Subscription{}, errors.ErrInvalidInput
		}
		end = &v
	}
	if end != nil && end.Before(start) {
		return model.Subscription{}, errors.ErrInvalidInput
	}

	sub := model.Subscription{
		ID:          uuid.NewString(),
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      in.UserID,
		StartDate:   start,
		EndDate:     end,
	}

	return s.repo.Create(ctx, sub)
}

func (s *Service) Get(ctx context.Context, in GetInput) (model.Subscription, error) {
	if in.ID == "" {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	return s.repo.Get(ctx, in.ID)
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (model.Subscription, error) {
	if in.ID == "" {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	if in.ServiceName == "" {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	if in.Price < 0 {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	if _, err := uuid.Parse(in.UserID); err != nil {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	start, err := parseMonthYear(in.StartDate)
	if err != nil {
		return model.Subscription{}, errors.ErrInvalidInput
	}
	var end *time.Time
	if in.EndDate != nil {
		v, err := parseMonthYear(*in.EndDate)
		if err != nil {
			return model.Subscription{}, errors.ErrInvalidInput
		}
		end = &v
	}
	if end != nil && end.Before(start) {
		return model.Subscription{}, errors.ErrInvalidInput
	}

	sub := model.Subscription{
		ID:          in.ID,
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      in.UserID,
		StartDate:   start,
		EndDate:     end,
	}

	return s.repo.Update(ctx, sub)
}

func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	if in.ID == "" {
		return errors.ErrInvalidInput
	}
	return s.repo.Delete(ctx, in.ID)
}

func (s *Service) List(ctx context.Context, in ListInput) ([]model.Subscription, error) {
	if in.Limit < 0 || in.Offset < 0 {
		return nil, errors.ErrInvalidInput
	}

	var fromT *time.Time
	if in.From != nil {
		v, err := parseMonthYear(*in.From)
		if err != nil {
			return nil, errors.ErrInvalidInput
		}
		fromT = &v
	}
	var toT *time.Time
	if in.To != nil {
		v, err := parseMonthYear(*in.To)
		if err != nil {
			return nil, errors.ErrInvalidInput
		}
		toT = &v
	}

	filter := ports.SubscriptionListFilter{
		UserID:      in.UserID,
		ServiceName: in.ServiceName,
		From:        fromT,
		To:          toT,
		Limit:       in.Limit,
		Offset:      in.Offset,
	}
	return s.repo.List(ctx, filter)
}

func (s *Service) Total(ctx context.Context, in TotalInput) (int64, error) {
	fromT, err := parseMonthYear(in.From)
	if err != nil {
		return 0, errors.ErrInvalidInput
	}
	toT, err := parseMonthYear(in.To)
	if err != nil {
		return 0, errors.ErrInvalidInput
	}
	if toT.Before(fromT) {
		return 0, errors.ErrInvalidInput
	}
	filter := ports.SubscriptionTotalFilter{
		UserID:      in.UserID,
		ServiceName: in.ServiceName,
		From:        fromT,
		To:          toT,
	}
	return s.repo.TotalCost(ctx, filter)
}

func parseMonthYear(s string) (time.Time, error) {
	if len(s) != 7 {
		return time.Time{}, fmt.Errorf("invalid")
	}
	if s[2] != '-' {
		return time.Time{}, fmt.Errorf("invalid")
	}
	month, err := time.Parse("01", s[0:2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid")
	}
	year, err := time.Parse("2006", s[3:7])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid")
	}
	return time.Date(year.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
