package cache

import (
	"errors"
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedRepository struct {
	Client *memcache.Client
	limit  int
	ttl    int
}

func NewMemcachedRepository(limit, ttl int, server ...string) *MemcachedRepository {
	memcachedClient := memcache.New(server...)

	if memcachedClient.Ping() != nil {
		log.Fatalf("Failed to connect to memcache server %s", server)
	}

	return &MemcachedRepository{
		Client: memcachedClient,
		limit:  limit,
		ttl:    ttl,
	}
}

func (r *MemcachedRepository) DecrementRequestCount(appKey string) (uint64, error) {
	val, err := r.Client.Decrement(appKey, 1)

	if errors.Is(err, memcache.ErrCacheMiss) {
		count := r.limit - 1

		err := r.Client.Set(
			&memcache.Item{
				Key:        appKey,
				Value:      []byte(fmt.Sprintf("%d", count)),
				Expiration: int32(r.ttl),
			})

		if err != nil {
			log.Fatalf("could not set value for key %s, error: %v\n", appKey, err)
			return 0, err
		}

		return uint64(r.limit), nil
	}

	if err != nil {
		return 0, err
	}

	return val, nil
}
