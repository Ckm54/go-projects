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

type mockPokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

func TestCommandCatch(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockStatusCode int
		mockResponse   mockPokemon
		initialPokedex map[string]Pokemon
		expectErr      string
		expectCaught   bool
		expectOutput   string
	}{
		{
			name:         "missing arguments",
			args:         []string{},
			expectErr:    "usage: catch <pokemon_name>",
			expectCaught: false,
		},
		{
			name: "pokemon already caught",
			args: []string{"pikachu"},
			initialPokedex: map[string]Pokemon{
				"pikachu": {Name: "pikachu", BaseExperience: 100},
			},
			expectOutput: "pikachu is already in your Pokedex!",
		},
		{
			name:           "successful fetch and possible catch",
			args:           []string{"charmander"},
			mockStatusCode: http.StatusOK,
			mockResponse:   mockPokemon{Name: "charmander", BaseExperience: 60},
			expectCaught:   true,
		},
		{
			name:           "API returns 404",
			args:           []string{"unknownmon"},
			mockStatusCode: http.StatusNotFound,
			expectErr:      "failed to fetch pokemon: 404 Not Found",
			expectCaught:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockStatusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			originalBaseURL := constants.BASEURL
			constants.BASEURL = server.URL
			defer func() { constants.BASEURL = originalBaseURL }()

			cfg := &Config{
				Cache:   pokecache.NewCache(2 * time.Second),
				Pokedex: make(map[string]Pokemon),
			}
			for k, v := range tt.initialPokedex {
				cfg.Pokedex[k] = v
			}

			err := commandCatch(cfg, tt.args)

			if tt.expectErr != "" {
				if err == nil || err.Error() != tt.expectErr {
					t.Fatalf("expected error %q, got %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectCaught {
				pName := tt.mockResponse.Name
				if _, ok := cfg.Pokedex[pName]; !ok {
					t.Errorf("expected %s to be in pokedex, but not found", pName)
				}
			}
		})
	}
}
