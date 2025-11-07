package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ckm54/go-projects/pokedexcli/constants"
	"github.com/ckm54/go-projects/pokedexcli/internal/pokecache"
)

func TestConfigUrl(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		action      Action
		expectedURL string
		expectedErr string
	}{
		{
			name:        "Next URL available",
			config:      &Config{Next: "https://example.com/next"},
			action:      ActionNext,
			expectedURL: "https://example.com/next",
		},
		{
			name:        "Reached end of list",
			config:      &Config{Next: "", Previous: "https://example.com/prev"},
			action:      ActionNext,
			expectedErr: "you have reached the end",
		},
		{
			name:        "Default start page when no Next or Previous",
			config:      &Config{Next: "", Previous: ""},
			action:      ActionNext,
			expectedURL: constants.BASEURL + "/location-area",
		},
		{
			name:        "Previous URL available",
			config:      &Config{Previous: "https://example.com/prev"},
			action:      ActionPrevious,
			expectedURL: "https://example.com/prev",
		},
		{
			name:        "First page, no Previous URL",
			config:      &Config{Previous: ""},
			action:      ActionPrevious,
			expectedErr: "you are on the first page",
		},
		{
			name:        "Invalid action returns error",
			config:      &Config{},
			action:      "unknown",
			expectedErr: "no action specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := configUrl(tt.config, tt.action)

			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if url != tt.expectedURL {
					t.Errorf("expected URL %q, got %q", tt.expectedURL, url)
				}
			}
		})
	}
}

// ---- Tests for fetchData
func TestFetchData(t *testing.T) {
	mockResponse := LocationAreaResponse{
		Next:     "http://example.com/next",
		Previous: "http://example.com/prev",
		Results: []LocationArea{
			{Name: "kanto"},
			{Name: "johto"},
		},
	}

	t.Run("Successful fetch updates config and caches result", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(mockResponse)
		}))
		defer server.Close()

		cfg := &Config{Cache: pokecache.NewCache(5 * time.Second)}
		err := fetchData(cfg, server.URL)
		if err != nil {
			t.Fatalf("fetchData returned error: %v", err)
		}

		if cfg.Next != mockResponse.Next {
			t.Errorf("expected Next %q, got %q", mockResponse.Next, cfg.Next)
		}
		if cfg.Previous != mockResponse.Previous {
			t.Errorf("expected Previous %q, got %q", mockResponse.Previous, cfg.Previous)
		}

		if _, ok := cfg.Cache.Get(server.URL); !ok {
			t.Errorf("expected response to be cached")
		}
	})

	t.Run("Cache hit avoids HTTP call", func(t *testing.T) {
		data, _ := json.Marshal(mockResponse)
		cache := pokecache.NewCache(5 * time.Second)
		cache.Add("http://cached-url", data)

		cfg := &Config{Cache: cache}
		err := fetchData(cfg, "http://cached-url")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Next != mockResponse.Next {
			t.Errorf("expected cached Next %q, got %q", mockResponse.Next, cfg.Next)
		}
	})

	t.Run("Server returns error status code (expected to fail)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "bad request", http.StatusBadRequest)
		}))
		defer server.Close()

		cfg := &Config{Cache: pokecache.NewCache(5 * time.Second)}
		err := fetchData(cfg, server.URL)
		if err == nil {
			t.Fatalf("expected error due to non-200 response, got nil")
		}
	})
}

func TestCommandMap(t *testing.T) {
	mockResponse := LocationAreaResponse{
		Next:     "http://example.com/next",
		Previous: "http://example.com/prev",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	cfg := &Config{Next: server.URL, Cache: pokecache.NewCache(5 * time.Second)}
	err := commandMap(cfg, []string{})
	if err != nil {
		t.Fatalf("commandMap returned error: %v", err)
	}

	if cfg.Next != mockResponse.Next {
		t.Errorf("expected Next %q, got %q", mockResponse.Next, cfg.Next)
	}
	if cfg.Previous != mockResponse.Previous {
		t.Errorf("expected Previous %q, got %q", mockResponse.Previous, cfg.Previous)
	}
}

func TestCommandMapB(t *testing.T) {
	mockResponse := LocationAreaResponse{
		Next:     "http://example.com/next2",
		Previous: "http://example.com/prev2",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	cfg := &Config{Previous: server.URL, Cache: pokecache.NewCache(5 * time.Second)}
	err := commandMapB(cfg, []string{})
	if err != nil {
		t.Fatalf("commandMapB returned error: %v", err)
	}

	if cfg.Next != mockResponse.Next {
		t.Errorf("expected Next %q, got %q", mockResponse.Next, cfg.Next)
	}
	if cfg.Previous != mockResponse.Previous {
		t.Errorf("expected Previous %q, got %q", mockResponse.Previous, cfg.Previous)
	}
}
