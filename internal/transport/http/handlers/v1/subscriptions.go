package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domainerrors "sub-hub/internal/domain/errors"
	"sub-hub/internal/domain/model"
	"sub-hub/internal/usecase/subscriptions"
)

type subscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type subscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type totalResponse struct {
	Total int64 `json:"total"`
}

func SubscriptionsRouter(log *zap.Logger, svc *subscriptions.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/total", func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		userID := optionalQuery(r, "user_id")
		serviceName := optionalQuery(r, "service_name")

		total, err := svc.Total(r.Context(), subscriptions.TotalInput{
			UserID:      userID,
			ServiceName: serviceName,
			From:        from,
			To:          to,
		})
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, totalResponse{Total: total})
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var req subscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
			return
		}
		created, err := svc.Create(r.Context(), subscriptions.CreateInput{
			ServiceName: req.ServiceName,
			Price:       req.Price,
			UserID:      req.UserID,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
		})
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, toResponse(created))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userID := optionalQuery(r, "user_id")
		serviceName := optionalQuery(r, "service_name")
		from := optionalQuery(r, "from")
		to := optionalQuery(r, "to")
		limit := parseIntQuery(r, "limit", 50)
		offset := parseIntQuery(r, "offset", 0)

		items, err := svc.List(r.Context(), subscriptions.ListInput{
			UserID:      userID,
			ServiceName: serviceName,
			From:        from,
			To:          to,
			Limit:       limit,
			Offset:      offset,
		})
		if err != nil {
			writeError(w, err)
			return
		}
		resp := make([]subscriptionResponse, 0, len(items))
		for _, it := range items {
			resp = append(resp, toResponse(it))
		}
		writeJSON(w, http.StatusOK, resp)
	})

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if _, err := uuid.Parse(id); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_id"})
				return
			}
			sub, err := svc.Get(r.Context(), subscriptions.GetInput{ID: id})
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, toResponse(sub))
		})

		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if _, err := uuid.Parse(id); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_id"})
				return
			}
			var req subscriptionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
				return
			}
			updated, err := svc.Update(r.Context(), subscriptions.UpdateInput{
				ID:          id,
				ServiceName: req.ServiceName,
				Price:       req.Price,
				UserID:      req.UserID,
				StartDate:   req.StartDate,
				EndDate:     req.EndDate,
			})
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, toResponse(updated))
		})

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if _, err := uuid.Parse(id); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_id"})
				return
			}
			if err := svc.Delete(r.Context(), subscriptions.DeleteInput{ID: id}); err != nil {
				writeError(w, err)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})

	_ = log
	return r
}

func toResponse(s model.Subscription) subscriptionResponse {
	start := formatMonthYear(s.StartDate)
	var end *string
	if s.EndDate != nil {
		v := formatMonthYear(*s.EndDate)
		end = &v
	}
	return subscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   start,
		EndDate:     end,
	}
}

func formatMonthYear(t time.Time) string {
	return t.Format("01-2006")
}

func optionalQuery(r *http.Request, key string) *string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	return &v
}

func parseIntQuery(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, domainerrors.ErrInvalidInput) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_input"})
		return
	}
	if errors.Is(err, domainerrors.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not_found"})
		return
	}
	if errors.Is(err, domainerrors.ErrConflict) {
		writeJSON(w, http.StatusConflict, errorResponse{Error: "conflict"})
		return
	}
	writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal"})
}
