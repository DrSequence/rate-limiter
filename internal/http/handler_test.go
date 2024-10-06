package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rate-limiter/domain"
	"rate-limiter/internal/app"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyHandler_Success(t *testing.T) {
	mockService := new(app.MockRateLimiterService)
	mockService.On("VerifyRequest", "12345", "my_app", "order_001").
		Return(true, nil)

	handler := NewRateLimiterHandler(mockService)

	reqBody := domain.RateLimitRequest{OrderID: "order_001"}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/verify", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("User-id", "12345")
	req.Header.Set("App-id", "my_app")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Verify(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.RateLimitResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Allowed)
	assert.Equal(t, "Request allowed", response.Message)

	mockService.AssertExpectations(t)
}

func TestVerifyHandler_LimitExceeded(t *testing.T) {
	mockService := new(app.MockRateLimiterService)
	mockService.On("VerifyRequest", "12345", "my_app", "order_001").
		Return(false, fmt.Errorf("request limit exceeded"))

	handler := NewRateLimiterHandler(mockService)

	reqBody := domain.RateLimitRequest{OrderID: "order_001"}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/verify", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("User-id", "12345")
	req.Header.Set("App-id", "my_app")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Verify(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	mockService.AssertExpectations(t)
}

func TestVerifyHandler_MissingHeaders(t *testing.T) {
	mockService := new(app.MockRateLimiterService)
	handler := NewRateLimiterHandler(mockService)

	reqBody := domain.RateLimitRequest{OrderID: "order_001"}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/verify", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Verify(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "user_id and app_id are required headers")
}
