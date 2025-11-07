package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommandPokedex(t *testing.T) {
	tests := []struct {
		name       string
		pokedex    map[string]Pokemon
		wantErr    string
		wantOutput string
	}{
		{
			name:       "empty pokedex",
			pokedex:    map[string]Pokemon{},
			wantErr:    "your pokedex is empty",
			wantOutput: "",
		},
		{
			name: "pokedex has one pokemon",
			pokedex: map[string]Pokemon{
				"pikachu": {Name: "pikachu"},
			},
			wantErr:    "",
			wantOutput: " -pikachu",
		},
		{
			name: "pokedex has multiple pokemons",
			pokedex: map[string]Pokemon{
				"charmander": {Name: "charmander"},
				"bulbasaur":  {Name: "bulbasaur"},
			},
			wantErr:    "",
			wantOutput: " -charmander -bulbasaur",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cfg := &Config{
				Pokedex: tt.pokedex,
				Out:     &out,
			}

			err := commandPokedex(cfg, nil)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotOutput := strings.TrimSpace(out.String())
			if tt.wantOutput != "" && !strings.Contains(gotOutput, strings.TrimSpace(tt.wantOutput)) {
				t.Errorf("expected output to contain %q, got %q", tt.wantOutput, gotOutput)
			}
		})
	}
}
