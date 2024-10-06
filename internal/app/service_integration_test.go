package app

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"rate-limiter/internal/cache"
	"testing"
)

func startMemcachedContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "memcached:alpine",
		ExposedPorts: []string{"11211/tcp"},
		WaitingFor:   wait.ForListeningPort("11211/tcp"),
	}

	memcachedC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	host, err := memcachedC.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := memcachedC.MappedPort(ctx, "11211")
	if err != nil {
		return nil, "", err
	}

	address := fmt.Sprintf("%s:%s", host, port.Port())
	return memcachedC, address, nil
}

func TestRateLimiterWithMemcached(t *testing.T) {
	ctx := context.Background()
	limit, ttl := 10, 60

	memcachedC, memcachedAddress, err := startMemcachedContainer(ctx)
	assert.NoError(t, err)

	defer func(memcachedC testcontainers.Container, ctx context.Context) {
		err := memcachedC.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(memcachedC, ctx)

	cacheRepo := cache.NewMemcachedRepository(limit, ttl, memcachedAddress)
	rateLimiterService := NewRateLimiterService(cacheRepo)

	tests := []struct {
		name      string
		userID    string
		appID     string
		orderID   string
		expectErr bool
		errMsg    string
		times     int
		isUnder   bool
	}{
		{
			name:    "First request, should be allowed",
			userID:  "12345",
			appID:   "my_app",
			orderID: "order_001",
			times:   1,
			isUnder: false,
		},
		{
			name:    "Multiple requests under limit",
			userID:  "12345",
			appID:   "my_app_10",
			orderID: "order_002",
			times:   10,
			isUnder: false,
		},
		{
			name:    "Request exceeds limit",
			userID:  "12345",
			appID:   "my_app_123",
			orderID: "order_10001",
			errMsg:  "request limit exceeded",
			times:   100,
			isUnder: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := 1

			// base case
			for index < tt.times && index < limit {
				allowed, testErr := rateLimiterService.
					VerifyRequest(tt.userID, tt.appID, fmt.Sprintf("order_%d", index+2))

				assert.NoError(t, testErr)
				assert.True(t, allowed)
				index++
			}

			// if more
			if tt.isUnder && index <= limit {
				allowed, testErr := rateLimiterService.
					VerifyRequest(tt.userID, tt.appID, fmt.Sprintf("order_%d", 9991))

				assert.NotNil(t, testErr)
				assert.False(t, allowed)
			}
		})
	}
}
