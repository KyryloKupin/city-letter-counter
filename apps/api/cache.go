package main

import (
	"context"
	"log"
	"slices"
	"sync"
	"time"
)

type cityGetter interface {
	getCities(ctx context.Context) ([]string, error)
}

type cityCache struct {
	source      cityGetter
	ttl         time.Duration
	retryAfter  time.Duration
	maxStaleAge time.Duration
	now         func() time.Time

	mu        sync.Mutex
	cities    []string
	fetchedAt time.Time
	expiresAt time.Time
}

func newCityCache(source cityGetter, ttl, retryAfter, maxStaleAge time.Duration) *cityCache {
	return &cityCache{
		source:      source,
		ttl:         ttl,
		retryAfter:  retryAfter,
		maxStaleAge: maxStaleAge,
		now:         time.Now,
	}
}

func (c *cityCache) getCities(ctx context.Context) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cities != nil && c.now().Before(c.expiresAt) {
		return slices.Clone(c.cities), nil
	}

	cities, err := c.source.getCities(ctx)
	if err != nil {
		if c.cities == nil {
			return nil, err
		}

		if c.now().Sub(c.fetchedAt) >= c.maxStaleAge {
			log.Print("Failed to refresh cities from GeoNames, the cached list is too old to serve")

			c.cities = nil

			return nil, err
		}

		log.Print("Failed to refresh cities from GeoNames, serving the cached list")

		c.expiresAt = c.now().Add(c.retryAfter)

		return slices.Clone(c.cities), nil
	}

	c.cities = slices.Clone(cities)
	c.fetchedAt = c.now()
	c.expiresAt = c.fetchedAt.Add(c.ttl)

	return cities, nil
}
