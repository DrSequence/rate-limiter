package domain

type RateLimitRequest struct {
	OrderID string `json:"order_id"`
}

type AppKey struct {
	UserID string
	AppID  string
}

type RateLimitResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message"`
	RetryIn int    `json:"retry_in,omitempty"`
}
