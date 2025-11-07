package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommandInspect(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		pokedex       map[string]Pokemon
		wantErr       bool
		expectedErr   string
		expectedPrint string
	}{
		{
			name:        "no arguments provided",
			args:        []string{},
			pokedex:     map[string]Pokemon{},
			wantErr:     true,
			expectedErr: "usage: catch <pokemon_name>",
		},
		{
			name: "pokemon not caught",
			args: []string{"pikachu"},
			pokedex: map[string]Pokemon{
				"charmander": {Name: "charmander"},
			},
			wantErr:     true,
			expectedErr: "you have not caught that pokemon",
		},
		{
			name: "pokemon caught and inspected",
			args: []string{"pikachu"},
			pokedex: map[string]Pokemon{
				"pikachu": {
					Name:   "pikachu",
					Height: 4,
					Weight: 60,
					Stats: []PokemonStat{
						{Stat: PokemonTypeInfo{Name: "speed"}, BaseStat: 90},
					},
					Types: []PokemonType{
						{Type: PokemonTypeInfo{Name: "electric"}},
					},
				},
			},
			wantErr:       false,
			expectedPrint: " -speed: 90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cfg := &Config{
				Pokedex: tt.pokedex,
				Out:     out,
			}

			err := commandInspect(cfg, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("expected error %q, got %q", tt.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(out.String(), tt.expectedPrint) {
					t.Errorf("expected output to contain %q, got %q", tt.expectedPrint, out.String())
				}
			}
		})
	}
}
