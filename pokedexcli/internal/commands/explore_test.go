package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ckm54/go-projects/pokedexcli/constants"
	"github.com/ckm54/go-projects/pokedexcli/internal/pokecache"
)

func TestCommandExplore(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		mockResponse any
		mockStatus   int
		preCache     bool
		preCacheData []byte
		expectErr    string
		expectCache  bool
	}{
		{
			name:        "missing arguments",
			args:        []string{},
			expectErr:   "usage: explore <area_name>",
			expectCache: false,
		},
		{
			name: "successful fetch and cache",
			args: []string{"kanto-route-1"},
			mockResponse: LocationAreaDetails{
				PokemonEncounters: []PokemonEncounter{
					{Pokemon: PokemonDetails{Name: "pikachu"}},
					{Pokemon: PokemonDetails{Name: "rattata"}},
				},
			},
			mockStatus:  http.StatusOK,
			expectCache: true,
		},
		{
			name:        "API returns error status",
			args:        []string{"unknown-area"},
			mockStatus:  http.StatusNotFound,
			expectErr:   "failed to fetch: 404 Not Found",
			expectCache: false,
		},
		{
			name:         "cache hit with valid JSON",
			args:         []string{"viridian-forest"},
			preCache:     true,
			preCacheData: []byte(`{"pokemon_encounters":[{"pokemon":{"name":"weedle"}}]}`),
			expectCache:  true,
		},
		{
			name:         "cache hit with invalid JSON",
			args:         []string{"broken-cache"},
			preCache:     true,
			preCacheData: []byte(`{invalid json`),
			expectErr:    "failed to decode cached response: invalid character 'i' looking for beginning of object key string",
			expectCache:  true,
		},
		{
			name:         "malformed API response",
			args:         []string{"bad-json"},
			mockResponse: `{"pokemon_encounters": [invalid]}`,
			mockStatus:   http.StatusOK,
			expectErr:    "invalid character 'i' looking for beginning of value",
			expectCache:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				switch v := tt.mockResponse.(type) {
				case string:
					fmt.Fprint(w, v)
				default:
					json.NewEncoder(w).Encode(v)
				}
			}))
			defer server.Close()

			originalBaseURL := constants.BASEURL
			constants.BASEURL = server.URL
			defer func() { constants.BASEURL = originalBaseURL }()

			cache := pokecache.NewCache(2 * time.Second)
			cfg := &Config{Cache: cache}

			if tt.preCache {
				url := fmt.Sprintf("%s/location-area/%s", constants.BASEURL, tt.args[0])
				cache.Add(url, tt.preCacheData)
			}

			err := commandExplore(cfg, tt.args)

			if tt.expectErr != "" {
				if err == nil || err.Error() != tt.expectErr {
					t.Fatalf("expected error %q, got %v", tt.expectErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tt.args) > 0 {
				url := fmt.Sprintf("%s/location-area/%s", constants.BASEURL, tt.args[0])
				_, found := cache.Get(url)
				if tt.expectCache && !found {
					t.Errorf("expected response to be cached, but not found")
				}
			}
		})
	}
}
