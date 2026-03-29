package subscriptions

type CreateInput struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
}

type UpdateInput struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
}

type GetInput struct {
	ID string
}

type DeleteInput struct {
	ID string
}

type ListInput struct {
	UserID      *string
	ServiceName *string
	From        *string
	To          *string
	Limit       int
	Offset      int
}

type TotalInput struct {
	UserID      *string
	ServiceName *string
	From        string
	To          string
}

type SubscriptionDTO struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}
