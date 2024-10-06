package app

import (
	"github.com/stretchr/testify/mock"
)

type MockRateLimiterService struct {
	mock.Mock
}

func (m *MockRateLimiterService) VerifyRequest(userID, appID, orderID string) (bool, error) {
	args := m.Called(userID, appID, orderID)
	return args.Bool(0), args.Error(1)
}
