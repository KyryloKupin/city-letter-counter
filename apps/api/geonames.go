package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type geoNamesClient struct {
	httpClient *http.Client
	endpoint   string
	username   string
}

func (c *geoNamesClient) getCities(ctx context.Context) ([]string, error) {
	endpoint, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse GeoNames endpoint: %w", err)
	}

	query := endpoint.Query()
	query.Set("featureClass", "P")
	query.Set("maxRows", "50")
	query.Set("orderby", "population")
	query.Set("username", c.username)

	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create GeoNames request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch GeoNames cities: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close GeoNames response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GeoNames returned HTTP %d", resp.StatusCode)
	}

	const maxResponseBytes = 1 << 20

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read GeoNames response: %w", err)
	}

	if len(body) > maxResponseBytes {
		return nil, fmt.Errorf("GeoNames response exceeds size limit")
	}

	var payload struct {
		GeoNames []struct {
			Name string `json:"name"`
		} `json:"geonames"`
		Status json.RawMessage `json:"status"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode GeoNames response: %w", err)
	}

	if payload.GeoNames == nil || (payload.Status != nil && string(payload.Status) != "null") {
		return nil, fmt.Errorf("GeoNames response contains an error or no city list")
	}

	cities := make([]string, 0, len(payload.GeoNames))

	for _, city := range payload.GeoNames {
		if city.Name == "" {
			return nil, fmt.Errorf("GeoNames response contains an empty city name")
		}

		cities = append(cities, city.Name)
	}

	return cities, nil
}
