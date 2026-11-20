package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"
)

const successCacheControl = "public, max-age=300"

type citiesResponse struct {
	StartsWith string   `json:"startsWith,omitempty"`
	Count      int      `json:"count"`
	Cities     []string `json:"cities"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func newHandler(source cityGetter) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/cities", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		startsWith := query.Get("startsWith")
		if query.Has("startsWith") && !isSingleLetter(startsWith) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "startsWith must be a single letter"})

			return
		}

		cities, err := source.getCities(r.Context())
		if err != nil {
			log.Print("Failed to fetch cities from GeoNames")

			writeJSON(w, http.StatusBadGateway, errorResponse{Error: "Failed to fetch cities"})

			return
		}

		if startsWith != "" {
			cities = filterCities(cities, startsWith)
		}

		writeJSON(w, http.StatusOK, citiesResponse{
			StartsWith: startsWith,
			Count:      len(cities),
			Cities:     cities,
		})
	})

	return mux
}

func isSingleLetter(value string) bool {
	letter, size := utf8.DecodeRuneInString(value)

	return size > 0 && size == len(value) && unicode.IsLetter(letter)
}

func filterCities(cities []string, letter string) []string {
	prefix := strings.ToLower(letter)
	matches := make([]string, 0)

	for _, city := range cities {
		if strings.HasPrefix(strings.ToLower(city), prefix) {
			matches = append(matches, city)
		}
	}

	return matches
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	if status == http.StatusOK {
		w.Header().Set("Cache-Control", successCacheControl)
	} else {
		w.Header().Set("Cache-Control", "no-store")
	}

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
