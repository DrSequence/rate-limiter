package http

import (
	"encoding/json"
	"net/http"

	"rate-limiter/domain"
	"rate-limiter/internal/app"
)

type RateLimiterHandler struct {
	rateLimiterService app.RateLimiterService
}

func NewRateLimiterHandler(service app.RateLimiterService) *RateLimiterHandler {
	return &RateLimiterHandler{rateLimiterService: service}
}

func (h *RateLimiterHandler) Verify(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-id")
	appID := r.Header.Get("App-id")

	if userID == "" || appID == "" {
		http.Error(w, "user_id and app_id are required headers", http.StatusBadRequest)
		return
	}

	var req domain.RateLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	allowed, err := h.rateLimiterService.VerifyRequest(userID, appID, req.OrderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	resp := domain.RateLimitResponse{
		Allowed: allowed,
		Message: "Request allowed",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
