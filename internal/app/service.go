package app

import (
	"fmt"
	"log"
	"time"

	"rate-limiter/internal/cache"
)

type RateLimiterService interface {
	VerifyRequest(userID, appID, orderID string) (bool, error)
}

type rateLimiterService struct {
	cacheRepo *cache.MemcachedRepository
}

func NewRateLimiterService(cacheRepo *cache.MemcachedRepository) RateLimiterService {
	return &rateLimiterService{cacheRepo: cacheRepo}
}

func (s *rateLimiterService) VerifyRequest(userID, appID, orderID string) (bool, error) {
	appKey := fmt.Sprintf("%s_%s", userID, appID)

	count, err := s.cacheRepo.DecrementRequestCount(appKey)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, fmt.Errorf("request limit exceeded")
	}

	return true, nil
}

func ResetCounters(cacheRepo *cache.MemcachedRepository) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			<-ticker.C
			err := cacheRepo.Client.FlushAll()
			if err != nil {
				log.Fatalln(err, "failed to flush cache")
				return
			} // flush counters
		}
	}()
}
